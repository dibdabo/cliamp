# Maintainer: bjarneo <https://github.com/bjarneo>
pkgname=cliamp
pkgver=1.63.2.r40.gdc5eaee
pkgrel=2
_commit=dc5eaee63d467c60845a7b8fe244f48ef7cdca53
_patch=0001-feat-search-add-multi-select-batch-actions.patch
pkgdesc='A retro terminal music player inspired by Winamp 2.x'
arch=('x86_64' 'aarch64')
url='https://github.com/bjarneo/cliamp'
license=('MIT')
depends=('alsa-lib' 'flac' 'libvorbis' 'libogg' 'ffmpeg' 'yt-dlp')
optdepends=('pipewire-alsa: audio output on PipeWire systems'
  'pulseaudio-alsa: audio output on PulseAudio systems')
makedepends=('go')
source=("${pkgname}-${_commit}.tar.gz::https://github.com/bjarneo/cliamp/archive/${_commit}.tar.gz"
        "${_patch}")
sha256sums=('c409cc009d69d54bc72fdbeded9f40afbc09839e23638e6fb275a4abcd053341'
            'fd438ade533f15fe7e6028d2708d0ea371552c0c36a61e6ddc6b4d2c55044484')

prepare() {
    cd "${pkgname}-${_commit}"
    patch -Np1 -i "${srcdir}/${_patch}"
}

build() {
    cd "${pkgname}-${_commit}"
    export CGO_ENABLED=1
    go build -trimpath -buildmode=pie \
        -ldflags="-s -w -X main.version=v${pkgver}.multiselect1 -linkmode=external" \
        -o cliamp .
}

package() {
    cd "${pkgname}-${_commit}"
    install -Dm755 cliamp "${pkgdir}/usr/bin/cliamp"
    install -Dm644 cliamp.desktop "${pkgdir}/usr/share/applications/cliamp.desktop"
    install -Dm644 Cliamp.png "${pkgdir}/usr/share/icons/hicolor/512x512/apps/cliamp.png"
    install -Dm644 Cliamp.png "${pkgdir}/usr/share/pixmaps/cliamp.png"
    install -Dm644 LICENSE "${pkgdir}/usr/share/licenses/${pkgname}/LICENSE"
}
