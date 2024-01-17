# Today

A command-line tool for working with the text of Christian scriptures. This tool
is closely associated with the website <openscripture.today>.

# Installation

Visit the [releases](https://github.com/zostay/today/releases) page for latest
binaries for your system.

Or you can install via the Golang compiler, if you have it installed:

```shell
go install github.com/zostay/today@latest
```

# Command-Line Usage

For basic help, just run the command:

```shell
today
```

That will output a menu displaying the available commands.

## Configuration

In order to use most commands, you will need to configure your [esv API token](https://api.esv.org/docs/) by placing it in the `ESV_API_TOKEN` environment variable or creating a file named `.esv.yaml` in your home directory containing your API key:

```yaml
access_token: YOUR_API_KEY
```

## Show a Verse

To display the content of a verse:

```shell
today show John 3:16
```

It will output the text. 

## Pick a Random Verse

To display a verse at random:

```shell
today random
```

This will display a random passage. You can use the `--book` option or the `--category` option to limit the random passage to a given book or category.

# Developer Tools

This project is a Golang library that provides a set of packages that may be used by other code to work with Bible references and Biblical text. This is a guide intended to introduce how these pieces work.

## Working with References

The typical entrypoint into referencing Biblical text is a standard verse reference. These references typically a form something like the following:

```
Luke 10:7; 1 Corinthians 9:4-11; Galations 6:6; 1 Timothy 5:17, 18
```

The `ref` package located in `github.com/zostay/today/pkg/ref` calls this a `ref.Multiple` reference and you can parse it using `ref.ParseMultiple`:

```go
refs, err := ref.ParseMultiple("Luke 10:7; 1 Corinthians 9:4-11; Galations 6:6; 1 Timothy 5:17,18")
if err != nil {
  panic(err)
}

// output the reconstructed
fmt.Println(refs.Ref())

// Validate() is called during parsing, but...
// validate that the reference appears sane
err := refs.Validate()
if err != nil {
  panic(err)
}

// Note: The above does not verify that the books named are books in any
// particular canon or refer to real verses, just that the references are
// formatted correctly.

// examine each part of the reference, for the above input this will output:
//
//   Luke 10:7
//   1 Corinthians 9:4-11
//   Galatians 6:6
//   1 Timothy 5:17,18
for _, r := range refs.Refs {
  fmt.Println(r.Ref())
}
```

References themselves do not have to refer to books in any particular canon and so validation merely states that the reference is plausibly correct. For example, `Luke 4:0` is an invalid reference, but `Philemon 12:4` and `Sterling 2:2` are both valid even though the first refers to a book without chapters and the second is a (perfectly valid) joke reference.

To turn the references into something more concrete, you can *resolve* the reference to a given canon. As of this writing, only the Protestant canon is available and currently defined in `ref.Canonical`. A canon lists all the books and valid verses for those books.

Therefore, if you want to take the reference parsed above and resolve it to complete references, you can do something like the following:

```go
res, err := ref.Canonical.Resolve(refs)
if err != nil {
  panic(err)
}

// examine each part of the reference, for the above input this will output:
//
//   Luke 10:7
//   1 Corinthians 9:4-9:11
//   Galatians 6:6
//   1 Timothy 5:17
//   1 Timothy 5:18
for _, r := range res {
	fmt.Println(r.Ref())
}
```

If there is no error during resolution, the named verses were all found within the canon.

If you want to understand the intricacies of how references are structured, see the Godoc reference.

## Biblical Text

Working with Biblical text does not require use of references. For that you can use the `text` package at `github.com/zostay/today/pkg/text`. As of this writing, this supports using the ESV API to retrieve Biblical text. To set up the ESV API, you will need to [get an API token](https://api.esv.org/docs/). You can either set this token in the `ESV_API_TOKEN` environment variable or create a file named `.esv.yaml` in your home directory, which contains your token like this:

```yaml
access_key: YOUR_API_KEY
```

To retrieve a specific reference, the basic code is as follows:

```go
package main

import (
	"fmt"
    
    "github.com/zostay/today/pkg/text/esv"
    "github.com/zostay/today/pkg/text"
)

func main() {
    res, err := esv.NewFromEnvironment()
    if err != nil {
        panic(err)
	}
    
    svc := text.NewService(res)
    txt, err := svc.Verse("John 3:16")
    if err != nil {
        panic(err)
	}
    
    fmt.Println(txt)
}
```

# Copyright & License

Copyright 2023 Andrew Sterling Hanenkamp.

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the “Software”), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
