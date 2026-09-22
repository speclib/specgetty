package ui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/mipmip/specgetty/src/scanner"
)

// --- 2.x a change arriving during a scan is not lost ---

// scanCount runs the update loop once and reports whether a scan was issued.
// A scan is a command, so its presence is what "a scan followed" means here.
func issuedScan(m model, msg tea.Msg) (model, bool) {
	updated, cmd := m.Update(msg)
	return updated.(model), cmd != nil
}

func watchingModel(t *testing.T) model {
	t.Helper()
	m := newModel(&scanner.Config{}, true, "0.0.0")
	m.width, m.height = 100, 30
	m.repoPaths = []string{"/p"}
	m.displayNames = []string{"p"}
	m.fields = append([]string(nil), defaultFields...)
	m.projects = scanner.ProjectMap{"/p": scanner.ProjectStatus{Info: scanner.ProjectInfo{
		Root: "/p", Origin: "/p",
	}}}
	m.recalcLayout()
	return m
}

func TestAChangeDuringAScanIsScannedAfterIt(t *testing.T) {
	m := watchingModel(t)

	m, started := issuedScan(m, fsChangeMsg{})
	if !started || !m.refreshing {
		t.Fatal("the first change should start a refresh")
	}
	if m.scanning {
		t.Fatal("a watcher rescan is a refresh, not a first scan")
	}

	// A second change while the first scan is still running. The watcher's
	// send is non-blocking into a channel of one, so without a flag this is
	// where the update is lost.
	m, again := issuedScan(m, fsChangeMsg{})
	if again && !m.scanPending {
		t.Fatal("a second scan should not start on top of the one in flight")
	}
	if !m.scanPending {
		t.Fatal("the change arriving during the scan should be remembered")
	}

	m, followed := issuedScan(m, scanMsg{projects: m.projects})
	if !followed {
		t.Error("a further scan should follow the one that landed")
	}
	if m.scanPending {
		t.Error("the pending change should be consumed by the scan it caused")
	}
	if !m.refreshing {
		t.Error("the further refresh is in flight")
	}
}

func TestManyChangesDuringAScanCauseOneFurtherScan(t *testing.T) {
	// A scan reads the whole project, so a queue of them would be a queue of
	// identical work.
	m := watchingModel(t)
	m, _ = issuedScan(m, fsChangeMsg{})
	for i := 0; i < 5; i++ {
		m, _ = issuedScan(m, fsChangeMsg{})
	}
	// One further scan, asserted positively so this cannot pass by never
	// remembering anything at all.
	m, followed := issuedScan(m, scanMsg{projects: m.projects})
	if !followed {
		t.Fatal("the burst should cause a further scan")
	}
	if m.scanPending {
		t.Error("five changes should leave one pending scan, not five")
	}
	m, more := issuedScan(m, scanMsg{projects: m.projects})
	if more {
		t.Error("the second landing should not start a third scan")
	}
	if m.scanning || m.refreshing {
		t.Error("nothing is in flight once the queue is drained")
	}
}

func TestAQuietProjectScansOnce(t *testing.T) {
	m := watchingModel(t)
	m, _ = issuedScan(m, fsChangeMsg{})
	m, followed := issuedScan(m, scanMsg{projects: m.projects})
	if followed {
		t.Error("nothing arrived during the scan, so nothing should follow it")
	}
	if m.scanning || m.refreshing {
		t.Error("the scan is done")
	}
}

// TestTheWatcherKeepsReadingWhileTheConsumerIsBusy is task 2.5. The send stays
// non-blocking: a watcher that waited on a busy consumer would stop reading
// filesystem events, and then it would miss changes rather than delay them.
func TestTheWatcherKeepsReadingWhileTheConsumerIsBusy(t *testing.T) {
	src, err := os.ReadFile(filepath.Join("..", "watcher", "watcher.go"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(src), "default:") {
		t.Error("the send into the events channel must stay non-blocking")
	}
}

// --- 3.x the store's git state on entering the properties tab ---

func TestGitStateIsReadOnEnteringThePropertiesTab(t *testing.T) {
	m := storeBackedModel()
	m.detailTab = tabProperties
	if cmds := m.refreshStoreGit(); len(cmds) != 1 {
		t.Fatalf("got %d reads, want one for the open store", len(cmds))
	}
	if cmds := m.enterTab(); len(cmds) == 0 {
		t.Error("entering the tab should read the store's git state")
	}
}

func TestNoGitStateIsReadForAProjectWithNoStore(t *testing.T) {
	m := plainModel()
	m.detailTab = tabProperties
	if cmds := m.refreshStoreGit(); len(cmds) != 0 {
		t.Errorf("got %d reads for a project holding its own content", len(cmds))
	}
}

func TestNothingIsReadBySessionThatNeverEntersTheTab(t *testing.T) {
	m := storeBackedModel()
	for _, tab := range []int{tabChanges, tabSpecs} {
		m.detailTab = tab
		if cmds := m.enterTab(); len(cmds) != 0 {
			t.Errorf("tab %d: got %d reads, want none until the tab is asked for", tab, len(cmds))
		}
	}
}

func TestAReadForAnotherProjectIsDropped(t *testing.T) {
	m := storeBackedModel()
	before := *m.projects[m.currentKey()].Info.Store
	m.applyStoreGit(storeGitMsg{project: "/somewhere/else",
		git: scanner.StoreGit{IsRepo: true, Ahead: 99}})
	after := *m.projects[m.currentKey()].Info.Store
	if after.Git != before.Git {
		t.Error("an answer for a project that is not open must not be shown against it")
	}
}

func TestAReadForTheOpenProjectIsRecorded(t *testing.T) {
	m := storeBackedModel()
	key := m.currentKey()
	m.applyStoreGit(storeGitMsg{project: key,
		git: scanner.StoreGit{IsRepo: true, TrackingKnown: true, Ahead: 3, Behind: 1}})

	g := m.projects[key].Info.Store.Git
	if g == nil || g.Ahead != 3 || g.Behind != 1 {
		t.Errorf("the re-read should be what the tab reports, got %+v", g)
	}
}

// --- 4.x the key is gone ---

func TestSDoesNothingOnEverySurface(t *testing.T) {
	for _, c := range []struct {
		name  string
		build func(t *testing.T) model
	}{
		{"change list", func(t *testing.T) model { return watchingModel(t) }},
		{"specs tab", func(t *testing.T) model {
			m := specDetailModel(t)
			return m
		}},
		{"a spec", func(t *testing.T) model {
			return press(specDetailModel(t), tea.KeyPressMsg{Code: tea.KeyEnter})
		}},
		{"a change's deltas", func(t *testing.T) model {
			return openedChangeSpecs(t, false)
		}},
	} {
		m := c.build(t)
		before := m
		after, cmd := m.Update(tea.KeyPressMsg{Code: 's', Text: "s"})
		got := after.(model)
		if cmd != nil {
			t.Errorf("%s: pressing s issued a command", c.name)
		}
		if got.scanning != before.scanning {
			t.Errorf("%s: pressing s started a scan", c.name)
		}
		if got.level != before.level || got.detailTab != before.detailTab {
			t.Errorf("%s: pressing s moved the view", c.name)
		}
	}
}

func TestTheNavBarNoLongerOffersScan(t *testing.T) {
	for _, c := range []struct {
		name  string
		build func(t *testing.T) model
	}{
		{"change list", func(t *testing.T) model { return watchingModel(t) }},
		{"specs tab", func(t *testing.T) model { return specDetailModel(t) }},
		{"a spec", func(t *testing.T) model {
			return press(specDetailModel(t), tea.KeyPressMsg{Code: tea.KeyEnter})
		}},
		{"a change's deltas", func(t *testing.T) model { return openedChangeSpecs(t, false) }},
	} {
		nav := ansi.Strip(c.build(t).renderNavBar())
		if strings.Contains(nav, "scan") {
			t.Errorf("%s: the nav bar still offers scan:\n%s", c.name, nav)
		}
	}
}

// TestTheRemainingRescanPathsStillWork is task 4.3. Removing the key must not
// remove the three callers that are not the key.
func TestTheRemainingRescanPathsStillWork(t *testing.T) {
	m := watchingModel(t)
	if cmd := m.rescanCurrent(); cmd == nil {
		t.Fatal("rescanCurrent should still produce a scan")
	}

	// The watcher.
	if _, started := issuedScan(m, fsChangeMsg{}); !started {
		t.Error("a filesystem change should still scan")
	}
	// An archive and a discard that succeeded.
	for _, msg := range []tea.Msg{
		archiveMsg{ok: true}, discardMsg{ok: true},
	} {
		if _, started := issuedScan(m, msg); !started {
			t.Errorf("%T should still scan", msg)
		}
	}
}

// TestThePickerIsStillTheManualPath is task 5.1. It re-resolves and re-reads
// the project it opens, including a store's git state, which is what `s` did.
func TestThePickerIsStillTheManualPath(t *testing.T) {
	src, err := os.ReadFile("picker.go")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(src), "m.doScanSingle(r.path)") {
		t.Error("opening a project from the picker must still read it")
	}
}
