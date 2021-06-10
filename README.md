RootDev Homepage
==================
Static website with a focus on least possible code in production. To lower
attack-vectors on the site and keeping maintenance as low as possible.

This means the website is 'pre-compiled' with Gulp (NodeJS).
* HTML is shrinked (removing spacing, comments);
* CSS is shrinked (removing spacing, comments and unused styles)
* All assets are GZIPped with Zopfli in advance (higher compression and no CPU time needed when served)

After that the site is offered through a custom self-written Go server.
* NGINx/OpenSSL, their both big projects in low-level languages with a lot
 of security issues lately;
* Less code == Less security issues, only have to watch the Go security announcements
 (and those are small and few);

Build steps
```
npm install -g yarn
cd static-src
yarn install
npm run build
cd ..
rm -rf pub
mv build pub
```
