package authentication

import (
	"fmt"
	"net/url"

	"github.com/hcarriz/reverb/authentication/provider"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/amazon"
	"github.com/markbates/goth/providers/apple"
	"github.com/markbates/goth/providers/auth0"
	"github.com/markbates/goth/providers/azuread"
	"github.com/markbates/goth/providers/battlenet"
	"github.com/markbates/goth/providers/bitbucket"
	"github.com/markbates/goth/providers/box"
	"github.com/markbates/goth/providers/dailymotion"
	"github.com/markbates/goth/providers/deezer"
	"github.com/markbates/goth/providers/digitalocean"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/dropbox"
	"github.com/markbates/goth/providers/eveonline"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/fitbit"
	"github.com/markbates/goth/providers/gitea"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gitlab"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/gplus"
	"github.com/markbates/goth/providers/heroku"
	"github.com/markbates/goth/providers/instagram"
	"github.com/markbates/goth/providers/intercom"
	"github.com/markbates/goth/providers/kakao"
	"github.com/markbates/goth/providers/lastfm"
	"github.com/markbates/goth/providers/line"
	"github.com/markbates/goth/providers/linkedin"
	"github.com/markbates/goth/providers/mastodon"
	"github.com/markbates/goth/providers/meetup"
	"github.com/markbates/goth/providers/microsoftonline"
	"github.com/markbates/goth/providers/naver"
	"github.com/markbates/goth/providers/nextcloud"
	"github.com/markbates/goth/providers/okta"
	"github.com/markbates/goth/providers/onedrive"
	"github.com/markbates/goth/providers/openidConnect"
	"github.com/markbates/goth/providers/patreon"
	"github.com/markbates/goth/providers/paypal"
	"github.com/markbates/goth/providers/salesforce"
	"github.com/markbates/goth/providers/seatalk"
	"github.com/markbates/goth/providers/shopify"
	"github.com/markbates/goth/providers/slack"
	"github.com/markbates/goth/providers/soundcloud"
	"github.com/markbates/goth/providers/spotify"
	"github.com/markbates/goth/providers/steam"
	"github.com/markbates/goth/providers/strava"
	"github.com/markbates/goth/providers/stripe"
	"github.com/markbates/goth/providers/tiktok"
	"github.com/markbates/goth/providers/twitch"
	"github.com/markbates/goth/providers/twitterv2"
	"github.com/markbates/goth/providers/typetalk"
	"github.com/markbates/goth/providers/uber"
	"github.com/markbates/goth/providers/vk"
	"github.com/markbates/goth/providers/wepay"
	"github.com/markbates/goth/providers/xero"
	"github.com/markbates/goth/providers/yahoo"
	"github.com/markbates/goth/providers/yammer"
	"github.com/markbates/goth/providers/yandex"
	"github.com/markbates/goth/providers/zoom"
	"golang.org/x/exp/slices"
)

// WithProvider adds a provider. Source is required for Okta, Nextcloud, and OpenID Providers.
func WithProvider(p provider.Provider, key, secret, callbackDomain, source string) Option {
	return option(func(a *auth) error {

		if !provider.Validate(p) {
			return ErrInvalidProvider
		}

		if slices.Contains([]provider.Provider{provider.Okta, provider.NextCloud, provider.OpenID}, p) {
			if _, err := url.Parse(source); err != nil {
				return err
			}
		}

		u, err := url.Parse(callbackDomain)
		if err != nil {
			return err
		}

		u.Path = fmt.Sprintf("/auth/callback/%s", p)

		d := cloneURL(u)
		d.Path = ""

		switch p {
		case provider.Amazon:
			goth.UseProviders(amazon.New(key, secret, u.String()))
		case provider.Apple:
			goth.UseProviders(apple.New(key, secret, u.String(), nil, apple.ScopeEmail, apple.ScopeName))
		case provider.Auth0:
			goth.UseProviders(auth0.New(key, secret, u.String(), d.String()))
		case provider.Azure:
			goth.UseProviders(azuread.New(key, secret, u.String(), nil))
		case provider.Battlenet:
			goth.UseProviders(battlenet.New(key, secret, u.String()))
		case provider.Bitbucket:
			goth.UseProviders(bitbucket.New(key, secret, u.String()))
		case provider.Box:
			goth.UseProviders(box.New(key, secret, u.String()))
		case provider.Dailymotion:
			goth.UseProviders(dailymotion.New(key, secret, u.String(), "email"))
		case provider.Deezer:
			goth.UseProviders(deezer.New(key, secret, u.String(), "email"))
		case provider.DigitalOcean:
			goth.UseProviders(digitalocean.New(key, secret, u.String(), "read"))
		case provider.Discord:
			goth.UseProviders(discord.New(key, secret, u.String(), discord.ScopeIdentify, discord.ScopeEmail))
		case provider.Dropbox:
			goth.UseProviders(dropbox.New(key, secret, u.String()))
		case provider.Eve:
			goth.UseProviders(eveonline.New(key, secret, u.String()))
		case provider.Facebook:
			goth.UseProviders(facebook.New(key, secret, u.String()))
		case provider.Fitbit:
			goth.UseProviders(fitbit.New(key, secret, u.String()))
		case provider.Gitea:
			goth.UseProviders(gitea.New(key, secret, u.String()))
		case provider.Github:
			goth.UseProviders(github.New(key, secret, u.String(), "read:user", "user:email"))
		case provider.Gitlab:
			goth.UseProviders(gitlab.New(key, secret, u.String()))
		case provider.Google:
			goth.UseProviders(google.New(key, secret, u.String()))
		case provider.GooglePlus:
			goth.UseProviders(gplus.New(key, secret, u.String()))
		case provider.Heroku:
			goth.UseProviders(heroku.New(key, secret, u.String()))
		case provider.Instagram:
			goth.UseProviders(instagram.New(key, secret, u.String()))
		case provider.Intercom:
			goth.UseProviders(intercom.New(key, secret, u.String()))
		case provider.Kakao:
			goth.UseProviders(kakao.New(key, secret, u.String()))
		case provider.LastFM:
			goth.UseProviders(lastfm.New(key, secret, u.String()))
		case provider.LINE:
			goth.UseProviders(line.New(key, secret, u.String(), "profile", "openid", "email"))
		case provider.Linkedin:
			goth.UseProviders(linkedin.New(key, secret, u.String()))
		case provider.Mastodon:
			goth.UseProviders(mastodon.New(key, secret, u.String(), "read:accounts"))
		case provider.Meetup:
			goth.UseProviders(meetup.New(key, secret, u.String()))
		case provider.Microsoft:
			goth.UseProviders(microsoftonline.New(key, secret, u.String()))
		case provider.Naver:
			goth.UseProviders(naver.New(key, secret, u.String()))
		case provider.NextCloud:
			goth.UseProviders(nextcloud.NewCustomisedDNS(key, secret, u.String(), source))
		case provider.Okta:
			goth.UseProviders(okta.New(key, secret, source, u.String()))
		case provider.Onedrive:
			goth.UseProviders(onedrive.New(key, secret, u.String()))
		case provider.OpenID:

			oic, err := openidConnect.New(key, secret, u.String(), source)
			if err != nil {
				return err
			}

			goth.UseProviders(oic)

		case provider.Patreon:
			goth.UseProviders(patreon.New(key, secret, u.String()))
		case provider.Paypal:
			goth.UseProviders(paypal.New(key, secret, u.String()))
		case provider.Salesforce:
			goth.UseProviders(salesforce.New(key, secret, u.String()))
		case provider.SeaTalk:
			goth.UseProviders(seatalk.New(key, secret, u.String()))
		case provider.Shopify:
			goth.UseProviders(shopify.New(key, secret, u.String()))
		case provider.Slack:
			goth.UseProviders(slack.New(key, secret, u.String()))
		case provider.SoundCloud:
			goth.UseProviders(soundcloud.New(key, secret, u.String()))
		case provider.Spotify:
			goth.UseProviders(spotify.New(key, secret, u.String()))
		case provider.Steam:
			goth.UseProviders(steam.New(key, u.String()))
		case provider.Strava:
			goth.UseProviders(strava.New(key, secret, u.String()))
		case provider.Stripe:
			goth.UseProviders(stripe.New(key, secret, u.String()))
		case provider.TikTok:
			goth.UseProviders(tiktok.New(key, secret, u.String()))
		case provider.Twitch:
			goth.UseProviders(twitch.New(key, secret, u.String()))
		case provider.Twitter:
			goth.UseProviders(twitterv2.New(key, secret, u.String()))
		case provider.Typetalk:
			goth.UseProviders(typetalk.New(key, secret, u.String(), "my"))
		case provider.Uber:
			goth.UseProviders(uber.New(key, secret, u.String()))
		case provider.VK:
			goth.UseProviders(vk.New(key, secret, u.String()))
		case provider.Wepay:
			goth.UseProviders(wepay.New(key, secret, u.String(), "view_user"))
		case provider.Xero:
			goth.UseProviders(xero.New(key, secret, u.String()))
		case provider.Yahoo:
			goth.UseProviders(yahoo.New(key, secret, u.String()))
		case provider.Yammer:
			goth.UseProviders(yammer.New(key, secret, u.String()))
		case provider.Yandex:
			goth.UseProviders(yandex.New(key, secret, u.String()))
		case provider.Zoom:
			goth.UseProviders(zoom.New(key, secret, u.String(), "read:user"))
		default:
			return ErrInvalidProvider
		}

		a.providers = append(a.providers, p)

		return nil
	})
}
