0.3.0  2024-01-22

 * Adding a `today openscripture photo` command with `--download`, `--on`, and `--yaml` flags for downloadinng and fetching metadata regarding the photo for the scripture of the day from [openscripture.today](https://openscripture.today).
 * Added the `TodayPhoto` command to the `ost` package for pulling photos from openscripture.today.
 * Added the `photo` and `unsplash` packages for working with photos, which are used to support the `ost` change and the CLI change, but also have future uses.

0.2.0  2024-01-19

 * Renaming `openscripture.today` to `openscripture` and adding `today` as a subcommand. (Both the upper level and sub-command do the same thing as the previous `openscripture.today` command.)
 * Adding an `opensripture on` command to fetch previous scriptures of the day from [openscripture.today](https://openscripture.today).
 * When parsing ranges, allow various unicode hyphens, not just U+002D.

0.1.0  2024-01-18

 * Added the `version` command to track which version is installed.
 * Added the `openscripture.today` (with `ost` alias) to allow for showing the scripture of the day from [openscripture.today](https://openscripture.today).
 * Added the API library in `pkg/ost` for working with the openscripture.today API.
 * Fix: There was a bug were resolution did not properly check to ensure that the verse was in the canon.

0.0.0  2024-01-17

 * Initial release.
 * Provides a library for working with Bible references in `pkg/ref`
 * Provides a tool for retrieving Biblical text from the ESV API
 * Provides a command-line tool named `today`
 * The `today` command can list books in the Protestant canon
 * The `today` command can list books by categories (of my selection)
 * The `today` command can retrieve a random passage from the ESV
 * The `today` command can show a specific passage from the ESV
