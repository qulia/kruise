package statefulsetrollout

import (
	"flag"
	stdlog "log"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/openkruise/kruise/pkg/apis"
	"github.com/openkruise/kruise/pkg/client/clientset/versioned/scheme"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	"github.com/onsi/gomega"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	one int32 = 1
	two int32 = 2
	ten int32 = 10
)

var cfg *rest.Config
var c client.Client

const timeout = time.Second * 2

// SetupTestReconcile returns a reconcile.Reconcile implementation that delegates to inner and
// writes the request to requests after Reconcile is finished.
func SetupTestReconcile(inner reconcile.Reconciler) (reconcile.Reconciler, chan reconcile.Request) {
	requests := make(chan reconcile.Request)
	fn := reconcile.Func(func(req reconcile.Request) (reconcile.Result, error) {
		result, err := inner.Reconcile(req)
		requests <- req
		return result, err
	})
	return fn, requests
}

// StartTestManager adds recFn
func StartTestManager(mgr manager.Manager, g *gomega.GomegaWithT) (chan struct{}, *sync.WaitGroup) {
	stop := make(chan struct{})
	wg := &sync.WaitGroup{}
	go func() {
		wg.Add(1)
		g.Expect(mgr.Start(stop)).NotTo(gomega.HaveOccurred())
		wg.Done()
	}()
	return stop, wg
}

func setUp(t *testing.T) (*gomega.GomegaWithT, chan reconcile.Request, chan struct{}, *sync.WaitGroup) {
	klog.InitFlags(nil)
	flag.Set("alsologtostderr", "true")
	flag.Set("v", "4")
	flag.Parse()
	g := gomega.NewGomegaWithT(t)
	// Setup the Manager and Controller.  Wrap the Controller Reconcile function so it writes each request to a
	// channel when it is finished.
	mgr, err := manager.New(cfg, manager.Options{})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	c = mgr.GetClient()
	recFn, requests := SetupTestReconcile(newReconciler(mgr))
	g.Expect(add(mgr, recFn)).NotTo(gomega.HaveOccurred())
	stopMgr, mgrStopped := StartTestManager(mgr, g)

	return g, requests, stopMgr, mgrStopped
}

func clean(g *gomega.GomegaWithT, c client.Client) {
	//
}

func TestMain(m *testing.M) {
	t := &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "..", "..", "config", "crds")},
	}
	apis.AddToScheme(scheme.Scheme)

	var err error
	if cfg, err = t.Start(); err != nil {
		stdlog.Fatal(err)
	}

	code := m.Run()
	t.Stop()
	os.Exit(code)
}
