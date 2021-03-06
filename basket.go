package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go"
	stripeClient "github.com/stripe/stripe-go/client"
	"github.com/tintinnabulate/aecontext-handlers/handlers"
	"github.com/tintinnabulate/gonfig"

	"golang.org/x/net/context"

	"google.golang.org/appengine/urlfetch"
)

// createHTTPRouter : create a HTTP router where each handler is wrapped by a given context
func createHTTPRouter(f handlers.ToHandlerHOF) *mux.Router {
	appRouter := mux.NewRouter()
	appRouter.HandleFunc("/", f(getHomePageHandler)).Methods("GET")
	appRouter.HandleFunc("/charge", f(postChargeHandler)).Methods("POST")
	return appRouter
}

// getHomePageHandler : show the homepage form
func getHomePageHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	tmpl := templates.Lookup("stripe.tmpl")
	tmpl.Execute(w, map[string]interface{}{
		"Key":            publishableKey,
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

// postChargeHandler : charge the customer
func postChargeHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("could not parse form: %v", err), http.StatusInternalServerError)
		return
	}

	emailAddress := r.Form.Get("stripeEmail")

	customerParams := &stripe.CustomerParams{Email: stripe.String(emailAddress)}
	customerParams.SetSource(r.Form.Get("stripeToken"))

	httpClient := urlfetch.Client(ctx)
	sc := stripeClient.New(stripe.Key, stripe.NewBackends(httpClient))

	newCustomer, err := sc.Customers.New(customerParams)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not create customer: %v", err), http.StatusInternalServerError)
		return
	}

	chargeParams := &stripe.ChargeParams{
		Amount:      stripe.Int64(int64(200)),
		Currency:    stripe.String("GBP"),
		Description: stripe.String("meeting payment"),
		Customer:    stripe.String(newCustomer.ID),
	}
	_, err = sc.Charges.New(chargeParams)
	if err != nil {
		fmt.Fprintf(w, "Could not process payment: %v", err)
		return
	}
	tmpl := templates.Lookup("payment_successful.tmpl")
	tmpl.Execute(w, map[string]interface{}{
		"Key":            publishableKey,
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

// Config is our configuration file format
type Config struct {
	CSRFKey              string `id:"CSRF_Key"             default:"my-random-32-bytes"`
	IsLiveSite           bool   `id:"IsLiveSite"           default:"false"`
	StripePublishableKey string `id:"StripePublishableKey" default:"pk_live_foo"`
	StripeSecretKey      string `id:"StripeSecretKey"      default:"sk_live_foo"`
	StripeTestPK         string `id:"StripeTestPK"         default:"pk_test_0n2vG3eX9wKGhKiB8hG0EhX2"`
	StripeTestSK         string `id:"StripeTestSK"         default:"rk_test_DXjSgrVJkA90FuBfu8NNf47H"`
}

var (
	publishableKey string
	templates      *template.Template
	config         Config
)

// configInit : initialise our config with the given config filename
func configInit(configName string) {
	err := gonfig.Load(&config, gonfig.Conf{
		FileDefaultFilename: configName,
		FileDecoder:         gonfig.DecoderJSON,
		FlagDisable:         true,
	})
	if err != nil {
		log.Fatalf("could not load configuration file: %v", err)
		return
	}
}

// routerInit : initialise our CSRF protected HTTPRouter
func routerInit() {
	router := createHTTPRouter(handlers.ToHTTPHandler)
	csrfProtector := csrf.Protect(
		[]byte(config.CSRFKey),
		csrf.Secure(config.IsLiveSite))
	csrfProtectedRouter := csrfProtector(router)
	http.Handle("/", csrfProtectedRouter)
}

// stripeInit : set up important Stripe variables
func stripeInit() {
	if config.IsLiveSite {
		publishableKey = config.StripePublishableKey
		stripe.Key = config.StripeSecretKey
	} else {
		publishableKey = config.StripeTestPK
		stripe.Key = config.StripeTestSK
	}
}

// templatesInit : parse the HTML templates, including any predefined functions (FuncMap)
func templatesInit() {
	templates = template.Must(template.New("").
		ParseGlob("templates/*.tmpl"))
}

func init() {
	configInit("config.json")
	templatesInit()
	routerInit()
	stripeInit()
}
