# Maintainer: Your Name <your@email.com>
pkgname=switchcalc
pkgver=1.0.0
pkgrel=1
pkgdesc="Complete calculator with Standard, Scientific, Programmer, and Date modes"
arch=('x86_64')
url="https://github.com/yourusername/switchcalc"
license=('MIT')
depends=('gtk4' 'glib2')
makedepends=('go')
source=()
sha256sums=()

build() {
    cd "$srcdir/.."
    go build -o switchcalc ./cmd/switchcalc
}

package() {
    cd "$srcdir/.."
    install -Dm755 switchcalc "$pkgdir/usr/bin/switchcalc"
    install -Dm644 switchcalc.desktop "$pkgdir/usr/share/applications/switchcalc.desktop"
}
