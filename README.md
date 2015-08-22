snips
-----------

A thing I built in a week of learning go.

Takes an input directory (current directory by default) finds all files ending in ".md"
runs them through the blackfriday markdown parser and saves the html in the output directory (current directory by default).  It also generates an index.html with links to all the markdown sorted by whatever go sorts files by in `ioutil.ReadDir`.

There's a good chance this code is awful.