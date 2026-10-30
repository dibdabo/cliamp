# Maintainer: jd
pkgname=cliamp
pkgver=2.3.0
pkgrel=4
pkgdesc='A retro terminal music player with local search and MPRIS customisations'
arch=('x86_64' 'aarch64')
url='https://github.com/bjarneo/cliamp'
license=('MIT')
depends=('alsa-lib' 'flac' 'libvorbis' 'libogg' 'mpg123' 'ffmpeg' 'yt-dlp')
optdepends=('pipewire-alsa: audio output on PipeWire systems'
            'pulseaudio-alsa: audio output on PulseAudio systems')
makedepends=('go')
source=("${pkgname}-${pkgver}.tar.gz::https://github.com/bjarneo/cliamp/archive/refs/tags/v${pkgver}.tar.gz"
        '0001-v2.3.0-preserve-customisations.patch')
sha256sums=('51828895ddb5b236fa0a119c1f06102c127065bf44911849c910fe798362fec5'
            '240dca632be44b81a2ef426533e11d240cefb67d79e12f950b5d2944c489f70a')

prepare() {
    cd "${pkgname}-${pkgver}"
    patch -Np1 -i "${srcdir}/0001-v2.3.0-preserve-customisations.patch"
}

build() {
    cd "${pkgname}-${pkgver}"
    export CGO_ENABLED=1
    go build -trimpath -buildmode=pie \
        -ldflags='-s -w -X main.version=v2.3.0.custom1 -linkmode=external' \
        -o cliamp .
}

package() {
    cd "${pkgname}-${pkgver}"
    install -Dm755 cliamp "${pkgdir}/usr/bin/cliamp"
    install -Dm644 cliamp.desktop "${pkgdir}/usr/share/applications/cliamp.desktop"
    install -Dm644 Cliamp.png "${pkgdir}/usr/share/icons/hicolor/512x512/apps/cliamp.png"
    install -Dm644 Cliamp.png "${pkgdir}/usr/share/pixmaps/cliamp.png"
    install -Dm644 LICENSE "${pkgdir}/usr/share/licenses/${pkgname}/LICENSE"
}
