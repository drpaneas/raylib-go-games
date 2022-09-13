# How to Contribute

We welcome contributions from the community. Here are a few ways you can help us improve.

## 📖 What changes we expect from you

Our first priority is to `transliterate` the code from [raylib-games] into [go-raylib].
Then write some tests (if possible) to ensure functionality remains the same and finally port the transilaterated code into idiomatic Go.

> Transliteration: It is originally a language term, meaning to write or print a letter or a word, using closest corresponding letters in a different alphabet or script.

### 1️⃣ Transliterate code from C to Go

It is a line-for-line (cosmetic only), compilation-errors-only-driven port of a piece of code from one language to another. Find the `main()` function and start working down from there.
This means you need to respect the original naming conventions and logic (**do not refactor**), but change only the absolute minimal code required to compile. If you diverge the code in any way that prevents like-to-like comparison, you are at disadvantage. This means you have to do all the *wrong* things, such as keep all the naming conventions of the original languages, the *case*-ing, and all the non-idiomatic qualities of the original codebase.

### 2️⃣ Transliterate tests (or write new ones)

The second step would be to port any tests or write new ones, to ensure that your code matches the expected functionality. The next step will be refactoring, so we need the safety net of the tests to make sure we do not diverge in functionallity. While doing so, feel free to add sufficient comments in the code.

### 3️⃣ Port into idiomatic Go

Refactor any code necessarry to look like **idiomatic** Go, by applying common Go [best-practices](https://github.com/golang/go/wiki/CodeReviewComments)
and [style](https://github.com/uber-go/guide/blob/master/style.md).

Make concepts easy to understand, readable and maintainanble.
We do care about the performance of the code, but don't go nuts on optimizing (e.g. no need to squeeze every single bit).
Have in mind this repository is meant for educational purposes, so the primary goal here is reability and usage of .

[raylib-games]: https://github.com/raysan5/raylib-games
[go-raylib]: https://github.com/gen2brain/raylib-go

---

## 📂 Open an Issue

If you see something you'd like changed, but aren't sure how to change it, submit an issue describing what you'd like to see.

Please create an issue with as much detail as you can provide. It should include the data gathered as indicated above, along with:

1. How to reproduce the problem
2. What the correct behavior should be
3. What the actual behavior is

We will do our very best to help you.

---

## 📝 Submit a Pull Request

If you find something you'd like to fix that's obviously broken, create a branch, commit your code, and submit a pull request.
Here's how:

1. Fork the repo on GitHub, and then clone it locally.
2. Create a branch named appropriately for the change you are going to make.
3. Make your code change.
4. If you are creating a function, please add a tests for it if possible.
5. Push your code change up to your forked repo.
6. Open a Pull Request to merge your changes to this repo. The comment box will be filled in automatically via a template.
7. All Pull Requests will be subject to Linting checks. Please make sure that your code complies and fix any warnings that arise. These are Checks that appear at the bottom of your Pull Request.
8. All Pull requests are subject to Testing. 

> Note, that we use `golangci-lint` to catch lint errors, and we require all contributors to install and use it.
Installation instructions can be found [here](https://golangci-lint.run/usage/install/).

See [Using Pull Requests](https://help.github.com/articles/using-pull-requests/) got more information on how to use GitHub PRs.

For an in depth guide on how to contribute see [this article](https://opensource.com/article/19/7/create-pull-request-github)