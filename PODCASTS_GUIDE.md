# cliamp podcast guide

This guide describes the dedicated **Podcasts** provider in the installed
cliamp v2.3.0. It does not require Spotify, an Apple account, an API key, or
a separate server.

## Start listening

Use the installed application:

```sh
cliamp --provider podcast
```

Or press uppercase `O` (`Shift+O`) in the player.

1. Press `/`, enter a podcast name, and press `Enter` to search.
2. Highlight a show and press `Enter` to load its episodes.
3. Highlight an episode and press `Enter` again to play.

Opening a show from the main provider list replaces the current playlist,
but does not immediately start playback. Press `Esc` or `b` to return to
the provider list.

## Discover shows

The provider offers a top chart and a **Browse Categories** entry. Discovery
uses Apple's directory; episode audio comes from publishers' RSS feeds.
Category browsing uses genre-name searches rather than separate genre charts.
Within category lists, `/` filters the entries already displayed.

The main provider's `/` search sends the query when you press `Enter` and
can return up to 100 shows. Press `Esc` to clear the search.

## Subscribe to shows

Highlight a show and press `f` to subscribe or unsubscribe. This works in
the provider list, category show lists, and `Ctrl+F` search results.
Saved shows appear under **Subscriptions**.

On an episode in the main playlist, `f` bookmarks that episode instead.

Subscriptions are stored in:

```text
~/.config/cliamp/podcast_subscriptions.json
```

You do not need to edit that file manually.

## Add an RSS feed directly

In the main Podcasts provider, press `/`, paste the publisher's RSS URL,
and press `Enter`. This bypasses Apple's directory and also works with
feed URLs that have no `.xml` or `.rss` suffix.

Press `f` on the resulting show to subscribe, or `Enter` to load episodes.
The feed must contain at least one playable episode.

You can also open a feed from the command line, or press `u` in the player
and paste its URL:

```sh
cliamp "https://example.com/podcast/feed.xml"
```

Use an actual RSS feed address, rather than a show's normal website or
Apple Podcasts listing.

## Configure startup and chart country

Edit `~/.config/cliamp/config.toml`. To open Podcasts by default and show
the UK chart:

```toml
# Top-level setting: put this before any [sections].
provider = "podcast"

[podcast]
country = "gb"
```

Merge these settings into existing sections rather than creating duplicates.
No `enabled` option is needed: the Podcasts provider is always available.

`country` is a two-letter code and defaults to `us`. It selects the chart
region only; it does not change searches, subscriptions, or publisher feeds.
cliamp does not automatically detect your country.

If `CLIAMP_CONFIG_DIR` or `XDG_CONFIG_HOME` is set, the config directory
may differ from the default paths used here.

## Search while building a queue

With Podcasts active, `Ctrl+F` searches for shows, with up to 20 results.
Unlike opening a show in the main provider list, these actions add episodes
to your existing playlist:

| Key on a show result | Action |
| --- | --- |
| `Enter` | Append episodes and start the first episode of that show |
| `a` | Append episodes; start the show's first episode if nothing is playing |
| `q` | Append episodes and queue them after existing queued tracks; start queued playback if idle |
| `f` | Subscribe or unsubscribe |
| `Esc` | Return to the query; press again to close |

## Playback and refresh controls

| Key | Action |
| --- | --- |
| `Space` | Pause or resume |
| `Enter` | Play the highlighted episode |
| `/` in the episode playlist | Search loaded episodes |
| `Left` / `Right` | Seek backward or forward 5 seconds |
| `Shift+Left` / `Shift+Right` | Seek backward or forward 30 seconds by default |
| `Ctrl+J` | Jump to a time |
| `[` / `]` | Decrease or increase playback speed |
| `a` on an episode | Toggle its play-next queue entry |
| `Ctrl+R` | Reload an open show, or refresh the chart in the provider list |

Reload applies to shows opened from the provider list or category browser.
For a show added through `Ctrl+F`, reopen its search result to fetch its feed
again; this adds episodes to the existing playlist.

## Limits and other podcast providers

The dedicated provider reads up to 300 playable episodes per feed, in the
publisher's order. Available titles, durations, episode numbers, and artwork
are retained. Items without playable audio are skipped.

Subscriptions persist, but this provider has no offline episode cache or
automatic episode downloads. Listening positions and played-state are saved
locally. Saved subscriptions can be listed without Apple's directory, but
fetching episodes and playing audio still requires access to the publisher.

For listening progress synced across devices, cliamp also supports a separate
Audiobookshelf provider (`B`). It requires an Audiobookshelf server and
credentials. Spotify (`S`) is another separate integration and is not needed
for the Podcasts provider described above.

## Listening progress and subscription shortcuts in v2.3.0

Episode positions are saved in `~/.config/cliamp/podcast_progress.json`.
Playing an unfinished episode again resumes five seconds before the saved
position. Stops within the first fifteen seconds restart from the beginning.
Finished episodes, or episodes stopped within their final minute, get a played
checkmark. Episodes under two minutes count as played after their midpoint.
Playing a completed episode again clears its played mark.

Press uppercase `F` to open subscriptions without replacing your current
playlist. In that overlay:

| Key | Action |
| --- | --- |
| `/` | Filter by show title or author |
| `Enter` | Append a show's episodes and start the first appended episode |
| `a` | Append all episodes of the highlighted show |
| `q` | Append and queue the highlighted show's episodes |
| `l` | Append and queue the highlighted show's newest episode |
| `L` | Append the newest episode from every subscription |

## This checkout and the installed application

This folder was previously on `spotify-search-multiselect`, an older branch
based on v1.63.2-era source. The installed application was independently
updated to v2.2.0. An intermediate `multiselect-mpris-v2.2.0` branch is also
preserved in Git.

The update branch `custom-v2.3.0` starts from upstream v2.3.0 and restores
search-result multi-selection, batch append/queue, batch-first help text,
the natural-language Spotify album hint, and the Linux media-control startup
fix. It keeps upstream's current album search implementation.

The old branches remain available for reference. `arch-package/PKGBUILD`
now builds v2.3.0 with a patch containing these custom changes. This folder
is needed for rebuilding or development, not for running the installed app.

## References

- [v2.3.0 podcast documentation](https://github.com/bjarneo/cliamp/blob/v2.3.0/docs/podcasts.md)
- [v2.3.0 podcast provider source](https://github.com/bjarneo/cliamp/blob/v2.3.0/external/podcast/provider.go)
- [v2.3.0 directory client source](https://github.com/bjarneo/cliamp/blob/v2.3.0/external/podcast/client.go)
