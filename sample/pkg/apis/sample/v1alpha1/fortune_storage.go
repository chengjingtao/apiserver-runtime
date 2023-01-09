package v1alpha1

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/apiserver/pkg/registry/rest"
	"net/http"
	"net/url"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource/resourcerest"
	"k8s.io/klog/v2"

)

// Rest functions
//
var _ resourcerest.Getter = &Fortune{}
var _ resourcerest.Lister = &Fortune{}
var _ resource.ObjectWithArbitrarySubResource = &Fortune{}
var _ resourcerest.TableConvertor = &Fortune{}

func  (f *Fortune)Destroy() {
}

func  (f *Fortune)New() runtime.Object {
	return &Fortune{}
}



func (in *Fortune) GetArbitrarySubResources() []resource.ArbitrarySubResource {
	return []resource.ArbitrarySubResource{
		&FortuneProxy{},
	}
}

// ConvertToTable handles table printing from kubectl get
func (f *Fortune) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	switch o := object.(type) {
	case *Fortune:
		return o.convertToTable(ctx, tableOptions)
	case *FortuneList:
		return o.convertToTable(ctx, tableOptions)
	}
	return nil, fmt.Errorf("unknown type Fortune %T", f)
}

// List implements rest.Lister
func (f *Fortune) List(ctx context.Context, o *internalversion.ListOptions) (runtime.Object, error) {
	//parts := strings.SplitN(o.LabelSelector.String(), "=", 2)

	//if len(parts) == 1 {
	//	fl := &FortuneList{}
	//	// return 5 random fortunes
	//	for i := 0; i < 5; i++ {
	//		obj, err := f.Get(ctx, "", &metav1.GetOptions{})
	//		if err != nil {
	//			return nil, err
	//		}
	//		fl.Items = append(fl.Items, *obj.(*Fortune))
	//	}
	//	return fl, nil
	//}
	//
	//fl := &FortuneList{}
	//var out []byte
	///* #nosec */
	//out, _ = exec.Command("/usr/games/fortune", "-s", "-m", parts[1]).Output()
	//values := strings.Split(string(out), "\n%\n")
	//for i, fo := range values {
	//	if i > 5 {
	//		break
	//	}
	//	if strings.TrimSpace(fo) == "" {
	//		continue
	//	}
	//	fl.Items = append(fl.Items, Fortune{Value: strings.TrimSpace(fo)})
	//}

	fl := &FortuneList{}
	return fl, nil
}

// Get implements rest.Getter
func (f *Fortune) Get(_ context.Context, name string, _ *metav1.GetOptions) (runtime.Object, error) {
	obj := &Fortune{}
	//var out []byte
	//// fortune exits non-zero on success
	//if name == "" {
	//	out, _ = exec.Command("/usr/games/fortune", "-s").Output()
	//} else {
	//	/* #nosec */
	//	out, _ = exec.Command("/usr/games/fortune", "-s", "-m", name).Output()
	//	fortunes := strings.Split(string(out), "\n%\n")
	//	if len(fortunes) > 0 {
	//		out = []byte(fortunes[0])
	//	}
	//}
	//if len(strings.TrimSpace(string(out))) == 0 {
	//	klog.Error("error out")
	//	return nil, errors.NewNotFound(Fortune{}.GetGroupVersionResource().GroupResource(), name)
	//}
	//obj.Value = strings.TrimSpace(string(out))
	obj.Value = "好好学习，天天向上"
	return obj, nil
}


// ---------- FortuneProxy

var _ resource.SubResource = &FortuneProxy{}
var _ rest.Storage = &FortuneProxy{}
var _ resourcerest.Connecter = &FortuneProxy{}

type FortuneProxy struct {
}


var proxyMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

// New returns an empty nodeProxyOptions object.
func (r *FortuneProxy) New() runtime.Object {
	return &FortuneProxyOptions{}
}

func  (f *FortuneProxy)Destroy() {
}

func (r *FortuneProxy) SubResourceName() string {
	return "proxy"
}

// ConnectMethods returns the list of HTTP methods that can be proxied
func (r *FortuneProxy) ConnectMethods() []string {
	return proxyMethods
}

// NewConnectOptions returns versioned resource that represents proxy parameters
func (r *FortuneProxy) NewConnectOptions() (runtime.Object, bool, string) {
	return &FortuneProxyOptions{}, true, "path"
}

// Connect returns a handler for the node proxy
func (r *FortuneProxy) Connect(ctx context.Context, id string, opts runtime.Object, responder rest.Responder) (http.Handler, error) {
	proxyOpts, ok := opts.(*FortuneProxyOptions)
	if !ok {
		return nil, fmt.Errorf("Invalid options object: %#v", opts)
	}
	//location, transport, err := node.ResourceLocation(r.Store, r.Connection, r.ProxyTransport, ctx, id)
	//if err != nil {
	//	return nil, err
	//}
	fmt.Printf("in connect function, id: %s, opts: %#v", id, opts)
	klog.Infof("in connect function, id: %s, opts: %#v", id, opts)
	location, _ := url.Parse("https://www.baidu.com/")
	location.Path = net.JoinPreservingTrailingSlash(location.Path, proxyOpts.Path)
	transport := http.DefaultTransport
	// Return a proxy handler that uses the desired transport, wrapped with additional proxy handling (to get URL rewriting, X-Forwarded-* headers, etc)
	return newThrottledUpgradeAwareProxyHandler(location, transport, true, false, responder), nil
}

func newThrottledUpgradeAwareProxyHandler(location *url.URL, transport http.RoundTripper, wrapTransport, upgradeRequired bool, responder rest.Responder) *proxy.UpgradeAwareHandler {
	handler := proxy.NewUpgradeAwareHandler(location, transport, wrapTransport, upgradeRequired, proxy.NewErrorResponder(responder))
	handler.MaxBytesPerSec = 0
	return handler
}
