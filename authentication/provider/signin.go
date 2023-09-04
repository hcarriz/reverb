package provider

import (
	"slices"
	"strings"

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
)

func (p Provider) Use(key, secret, url, source string, scope ...string) error {

	for x := range scope {
		scope[x] = strings.ToLower(scope[x])
	}

	switch p {
	case Amazon:
		goth.UseProviders(amazon.New(key, secret, url, scope...))
	case Apple:
		scope = append(scope, apple.ScopeEmail, apple.ScopeName)
		slices.Compact(scope)
		goth.UseProviders(apple.New(key, secret, url, nil, scope...))
	case Auth0:
		goth.UseProviders(auth0.New(key, secret, url, source, scope...))
	case Azure:
		goth.UseProviders(azuread.New(key, secret, url, nil, scope...))
	case Battlenet:
		goth.UseProviders(battlenet.New(key, secret, url, scope...))
	case Bitbucket:
		goth.UseProviders(bitbucket.New(key, secret, url, scope...))
	case Box:
		goth.UseProviders(box.New(key, secret, url, scope...))
	case Dailymotion:
		scope = append(scope, "email")
		slices.Compact(scope)
		goth.UseProviders(dailymotion.New(key, secret, url, scope...))
	case Deezer:
		scope = append(scope, "email")
		slices.Compact(scope)
		goth.UseProviders(deezer.New(key, secret, url, scope...))
	case DigitalOcean:
		scope = append(scope, "read")
		slices.Compact(scope)
		goth.UseProviders(digitalocean.New(key, secret, url, scope...))
	case Discord:
		scope = append(scope, discord.ScopeBot, discord.ScopeConnections, discord.ScopeEmail)
		slices.Compact(scope)
		goth.UseProviders(discord.New(key, secret, url, scope...))
	case Dropbox:
		goth.UseProviders(dropbox.New(key, secret, url, scope...))
	case Eve:
		goth.UseProviders(eveonline.New(key, secret, url, scope...))
	case Facebook:
		goth.UseProviders(facebook.New(key, secret, url, scope...))
	case Fitbit:
		goth.UseProviders(fitbit.New(key, secret, url, scope...))
	case Gitea:
		goth.UseProviders(gitea.New(key, secret, url, scope...))
	case Github:
		scope = append(scope, "read:user", "user:email")
		slices.Compact(scope)
		goth.UseProviders(github.New(key, secret, url, scope...))
	case Gitlab:
		goth.UseProviders(gitlab.New(key, secret, url, scope...))
	case Google:
		goth.UseProviders(google.New(key, secret, url, scope...))
	case GooglePlus:
		goth.UseProviders(gplus.New(key, secret, url, scope...))
	case Heroku:
		goth.UseProviders(heroku.New(key, secret, url, scope...))
	case Instagram:
		goth.UseProviders(instagram.New(key, secret, url, scope...))
	case Intercom:
		goth.UseProviders(intercom.New(key, secret, url, scope...))
	case Kakao:
		goth.UseProviders(kakao.New(key, secret, url, scope...))
	case LastFM:
		goth.UseProviders(lastfm.New(key, secret, url))
	case LINE:
		scope = append(scope, "profile", "openid", "email")
		slices.Compact(scope)
		goth.UseProviders(line.New(key, secret, url, scope...))
	case Linkedin:
		goth.UseProviders(linkedin.New(key, secret, url, scope...))
	case Mastodon:
		scope = append(scope, "read:accounts")
		slices.Compact(scope)
		goth.UseProviders(mastodon.New(key, secret, url, scope...))
	case Meetup:
		goth.UseProviders(meetup.New(key, secret, url, scope...))
	case Microsoft:
		goth.UseProviders(microsoftonline.New(key, secret, url, scope...))
	case Naver:
		goth.UseProviders(naver.New(key, secret, url))
	case NextCloud:
		goth.UseProviders(nextcloud.NewCustomisedDNS(key, secret, url, source))
	case Okta:
		goth.UseProviders(okta.New(key, secret, source, url, scope...))
	case Onedrive:
		goth.UseProviders(onedrive.New(key, secret, url, scope...))
	case OpenID:
		oic, err := openidConnect.New(key, secret, url, source, scope...)
		if err != nil {
			return err
		}
		goth.UseProviders(oic)
	case Patreon:
		goth.UseProviders(patreon.New(key, secret, url, scope...))
	case Paypal:
		goth.UseProviders(paypal.New(key, secret, url, scope...))
	case Salesforce:
		goth.UseProviders(salesforce.New(key, secret, url, scope...))
	case SeaTalk:
		goth.UseProviders(seatalk.New(key, secret, url, scope...))
	case Shopify:
		goth.UseProviders(shopify.New(key, secret, url, scope...))
	case Slack:
		goth.UseProviders(slack.New(key, secret, url, scope...))
	case SoundCloud:
		goth.UseProviders(soundcloud.New(key, secret, url, scope...))
	case Spotify:
		goth.UseProviders(spotify.New(key, secret, url, scope...))
	case Steam:
		goth.UseProviders(steam.New(key, url))
	case Strava:
		goth.UseProviders(strava.New(key, secret, url, scope...))
	case Stripe:
		goth.UseProviders(stripe.New(key, secret, url, scope...))
	case TikTok:
		goth.UseProviders(tiktok.New(key, secret, url, scope...))
	case Twitch:
		goth.UseProviders(twitch.New(key, secret, url, scope...))
	case Twitter:
		goth.UseProviders(twitterv2.New(key, secret, url))
	case Typetalk:
		scope = append(scope, "my")
		slices.Compact(scope)
		goth.UseProviders(typetalk.New(key, secret, url, scope...))
	case Uber:
		goth.UseProviders(uber.New(key, secret, url, scope...))
	case VK:
		goth.UseProviders(vk.New(key, secret, url, scope...))
	case Wepay:
		scope = append(scope, "view_user")
		slices.Compact(scope)
		goth.UseProviders(wepay.New(key, secret, url, scope...))
	case Xero:
		goth.UseProviders(xero.New(key, secret, url))
	case Yahoo:
		goth.UseProviders(yahoo.New(key, secret, url, scope...))
	case Yammer:
		goth.UseProviders(yammer.New(key, secret, url, scope...))
	case Yandex:
		goth.UseProviders(yandex.New(key, secret, url, scope...))
	case Zoom:
		scope = append(scope, "read:accounts")
		slices.Compact(scope)
		goth.UseProviders(zoom.New(key, secret, url, scope...))
	}

	return nil

}
