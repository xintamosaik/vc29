# readme

## plans

- [x] SPA
- [x] templ as a react replacement
- [x] go backend
- [x] htmx to extend html
- [x] tailwind
- [ ] maybe typescript or jsdoc as a ts replacement
- [ ] maybe air for hot reloading

## ideas and concepts

### above the fold

Websites should load fast. UX and SEO are the predominant reasons for that. There are several ways to get that. Trimming the initial bundle helps, but snappy UX demands we also preload or fetch the rest fast—ideally before the user notices.

So we lazily load scripts and styles for those parts of the application that will reveal themselves only after interaction and and we hope users take just long enough for lazy-loaded assets to finish before they scroll or click.

### lazy loading

Lazy loading has its quirks. A script tag that gets activated during initial load will block rendering per default. If we want to prevent that we have to add instructions to their script tags and for CSS for the link tags that invoke style files.

For scripts we can add "defer" or "async" and that solves things but it's good to know what these attributes mean and when to use which.

> Defer delays script execution until the DOM is parsed. async runs the script as soon as it’s fetched—useful for analytics, but risky for DOM-dependent code.

For styles there are a couple of ways to game lighthouse metrics but to make styles non-blocking we need to think ourselves. I know.

### prefetch-ing


## workarounds (unplanned features)

- [x] esbuild to bundle js and css

## setup

You can actually skip this step, as Go will automatically download dependencies when you run the code.
However, if you want to download them manually, you can run:

```bash
go mod download
```

## run

You can run the code with and get a development server on port 3000:

```bash
go tool templ generate && go run .
```

Also for now you need to run this for tailwind:

```bash
npx @tailwindcss/cli -i ./src.css -o ./dist/under_the_fold.css --watch
```

## build

I never tried it actually.

```bash
go tool templ generate && go build -o bin/app
```
