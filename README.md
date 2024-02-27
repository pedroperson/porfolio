# template-go-vercel

## Compile html

The /public folder is being used to serve all the static files.

To simplify the deployment, we want to compile all the go templates to html ahead of time

```
go run compiler/compiler.go
```

TODO: Write a watcher for this go file

## Compile css

Currently we are using tailwind, that compiles only on the compiled index.html

```
yarn tailwindcss -i templates/main.css -o public/style.css --watch
```

TODO: When compiling to deploy use:

```
yarn tailwindcss -i templates/main.css -o public/style.css --watch
```

## Run local server

```
yarn dev
```

## So watchu gotta do rn :

Open a terminal at root and run to start the tailwind listener

```
yarn tailwindcss -i templates/main.css -o public/style.css --watch
```

Then open another and run to open the local server

```
yarn dev
```

Then open another and run every time you change an html file or go:

```
go run compiler/compiler.go
```

## Publishing

Just push to main and vercel will do the rest

## Connecting domain to vercel

https://medium.com/@aozora-med/how-to-set-up-namecheap-domain-on-vercel-2b2313e22342
