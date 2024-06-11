## WIP  TBD

 * :computer: The `today ost index` command has been added to pull down scripture indexes from openscripture.today.
 * :computer: Added `--exclude` and `--exclude-index` options to `today random`.
 * ESV text from VerseText and others won't include the references now.
 * Updating the Verse and Photo file formats used by the `ost.Client` for the upcoming version of openscripture.today. Adding a `ost.Metadata` structure that records version, original verseion, and whether the verse and photo have been pruned.
 * A new `VersesIndex` has been added to the `ost.Client` for fetching indexes listing all the verses that have been postd for certain periods (all time, year, month).
 * Added a `Clone` method to `ref.Canon` and `ref.Book` to create deep clones of these objects.
 * Added the `ref.ExcludeReferences` option for use with the `ref.Random` functions.
 * Added a `Subtract` method to `ref.Resolved` to return the difference between two resolved references.
 * Added a `Filtered` method to `ref.Canon` that returns a new canon that has segments references removed.
 * :boom: Breaking Change :boom:: Now requires Go 1.22.
 * :boom: Breaking Change :boom:: Significant changes have been made to the photos API. These changes include the following:
   - The `photo.Meta` and `photo.Info` structures have been removed/merged/refactored into `photo.Descriptor`.
   - A new interface named `photo.Image` has been added and each `photo.Descriptor` should have at least one, but may have many associated `photo.Image` objects.
   - Interaction with the images associated with a `photo.Descriptor` is handled through the methods named `AddImage`, `RemoveImage`, `HasImage`, and `GetImage`.
   - Every `photo.Iamge` implementation must also implement either or both of `photo.ImageReader` or `photo.ImageDecoded`.
   - When stored with a `photo.Descriptor`, the `photo.Image` is transformed into a `photo.ImageComplete` if it is not one already.
   - The `photo.Iamge` interface required a `Filename` method to be implemented.
   - The `photo.ImageReader` interface requires a `Reader` method to be implemented.
   - The `photo.ImageDecoded` interface requires a `Image` method to be implemented.
   - The `photo.ImageComplete` interface is a union of `photo.ImageReader` and `photo.ImageDecoded`.
   - The `photo.CompleteImage` and associated `photo.Complete` function are provided to transform any `photo.Image` into a `photo.ImageComplete` (or return it as is if it already is such). 
   - The `photo.Service` no longer requires a `CacheKey` method to be implemented.
   - The `photo.Service` no longer requires a `Download` method to be implemented as that funcationlity should be handled via `photo.Image`.
   - The `photo.DominantImageColor` function has been added, replacing the removed `DominantImageColor` method of `photo.Service`.
   - The `ost.Client` returns `photo.Descriptor` objects instead of `photo.Info` objects.
   - Mirroring the API in Golang's built-in `image` package, a `photo.RegisterEncoder` function has been added to allow for automatic encoding. This includes the "jpeg" encoding by default and the `photo.DefaultEncoding` is used in situations when no encoding has been pre-selected, which is set to "jpeg" format.
   - A general purpose implementation of `photo.Image`/`photo.ImageDecoded` is provided in `photo.Memory`.
   - A general purpose implementation of `photo.Image`/`photo.ImageReader` is provided in `photo.File`.
   - Every `photo.Descriptor` should define an image for the `photo.Original` key, which should provide a full-size, original image.
   - The `photo.FromImage` option has been added for use with the `ResizeImage` method of `photo.Service`, which selects which `photo.Image` associated with a `photo.Descriptor` should be used as the source image for resizing.
 * New utility function `unsplash.IDFromURL()` added which will give you the photo ID from an Unsplash photo URL.
 * :hammer: Fix: Multiple total chapter references will be resolved to chapter ranges. For example, if you parse "Ps. 12-13", the resolver will correctly return "Psalms 12-13" instead of "Psalms 12:1-13:6" as it would have before.
 * :hammer: Fix: Psalms singular handling has been special cased so that references to a single verse can show something like "Psalm 12" or "Psalm 12:1-3" rather than "Psalms 12" or "Psalms 12:1-3" as it would have before.
 * :hammer: Fix: Photo descriptors were not being correctly written to JSON files prior to this version. They were missing the "creator" key and embedding the name of the creator in the structure above. This has been corrected.
 * :hammer: Fix: Verse JSON files were not written correctly either prior to this version. They were missing the "version" key and embedding the name of the version in the parent structure. This has been fixed.

## 0.5.1  2024-02-20

 * :hammer: Fix: The ESV resolver returns cleaner verse references now.

## 0.5.0  2024-02-20

 * :computer: The `today show` command supports most common book abbreviations now.
 * The `text.Service` allows for the canon and the abbreviations used to be configured using the `text.WithCanon`, `text.WithAbbreviations`, and `text.WithoutAbbreviations` options. The service still uses `ref.Canonical` by default and now uses `ref.Abbreviations` by default.
 * Addeded a `CompactRef()` method to `ref.Resolved` to return a compact string representation of the reference. (For example, Genesis 12:4-12:6 would be Genesis 12:4-6 or Genesis 12:1-20 would be Genesis 12.)
 * Addeda a `LastVerseInChapter()` method to `ref.Book` to return the last verse in a chapter (or in the case of a chapterless book like Obadiah or Philemon, the last verse in the book).
 * Added the `ref.BookAbbreviations` structure with associated components and the `PreferredAbbreviation()` and `BookName()` methods to assist with abbreviating and parsing abbreviated book names.
 * Added a standard set of abbreviations in `ref.Abbreviations`.
 * Added options to the `Resolve()` method of `ref.Canon`, the first (and only, for now) option is `ref.WithAbbreviations()` to select another set of abbreviations (rather than using ref.Abbreviations as is the default) or to use no abbreviations at all via `ref.WithoutAbbreviations()`.
 * Added the resolve options to the `Book()` method of `ref.Canon` as well.
 * Added a low-level interface called `ref.AbbrTree` for quickly resolving abbreviated book names.
 * :question::boom: Potentially Breaking Change: Implementations of `ref.Absolute` now must implement the `FullNameRef()` and `AbbreviatedRef()` methods.

## 0.4.0  2024-02-01

 * :boom: Breaking Change :boom:: The `ref.RandomPassage` ane `ref.RandomPassageFromRef` functions now take two additional integer arguments to select width of range returned.
 * :boom: Breaking Change :boom:: The `ost.Verse` structure has been moved to `text.Verse`. It has also been restructured to include a `Link` field and the HTML `Content` field has been replaced with a structure that contains `Text` and `HTML` fields.
 * :boom: Breaking Change :boom:: The `text.Service` has changed substantially: 
   - Methods `Verse`, `VerseText`, `VerseHTML`, `RandomVerse`, `RandomVerseText`, and `RandomVerseHTML` now require a `context.Context` argument.
   - The `Verse` method has been renamed to `VerseText` and a new `Verse` method that returns `text.Verse` has been added.
   - The `RandomVerse` method has been renamed to `RandomVerseText` and a new `RandomVerse` method that returns `text.Verse` has been added.
 * :boom: Breaking Change :boom:: The `text.Resolver` interface has changed substantially:
   - Requires a `context.Context` argument for all methods.
   - Requires a `VersionInformation` method.
   - Renamed `Verse` to `VerseText` and a new `Verse` method that returns `text.Verse` has been added.
 * :boom: Breaking Change :boom:: The `esv.Resolver` implements new `text.Resolver` changes.
 * :boom: Breaking Change :boom:: The `ost.Version` structure has been moved to `text.Version`.
 * :boom: Breaking Change :boom:: The `ost.Client` methods `Today`, `TodayVerse`, `TodayHTML`, and `TodayPhoto` now require a `context.Context` argument.
 * Added the `--minimum-verses` and `--maximum-verses` options to `today random` to allow control over how many verses are selected for the passage.
 * Added `ref.WithAtLeast()` and `ref.WithAtMost()` options to `ref.Random` to allow control over how many verses are selected for the passage.
 * Added a field for loading and saving the preferred color of an image to `photo.Meta` and the `SetColor` method for converting `color.Color` to a CSS hex color code and `GetColor` to perform the same operation in reverse.
 * When checking for dominant color in a photo, black and white are generally disqualified.
 * :hammer: Fix: Fixed a bug where the `today random` output showed the passage reference twice.
 * :hammer: Fix: Random passages should now be unbiased (previously, there was a slight bias towards picking passages at the end of a book or passage).

## 0.3.0  2024-01-22

 * :computer: Adding a `today openscripture photo` command with `--download`, `--on`, and `--yaml` flags for downloadinng and fetching metadata regarding the photo for the scripture of the day from [openscripture.today](https://openscripture.today).
 * Added the `TodayPhoto` command to the `ost` package for pulling photos from openscripture.today.
 * Added the `photo` and `unsplash` packages for working with photos, which are used to support the `ost` change and the CLI change, but also have future uses.

## 0.2.0  2024-01-19

 * :computer: Renaming `openscripture.today` to `openscripture` and adding `today` as a subcommand. (Both the upper level and sub-command do the same thing as the previous `openscripture.today` command.)
 * :computer: Adding an `opensripture on` command to fetch previous scriptures of the day from [openscripture.today](https://openscripture.today).
 * When parsing ranges, allow various unicode hyphens, not just U+002D.

## 0.1.0  2024-01-18

 * :computer: Added the `version` command to track which version is installed.
 * :computer: Added the `openscripture.today` (with `ost` alias) to allow for showing the scripture of the day from [openscripture.today](https://openscripture.today).
 * Added the API library in `pkg/ost` for working with the openscripture.today API.
 * :hammer: Fix: There was a bug were resolution did not properly check to ensure that the verse was in the canon.

## 0.0.0  2024-01-17

 * Initial release.
 * Provides a library for working with Bible references in `pkg/ref`
 * Provides a tool for retrieving Biblical text from the ESV API
 * :computer: Provides a command-line tool named `today`
 * :computer: The `today` command can list books in the Protestant canon
 * :computer: The `today` command can list books by categories (of my selection)
 * :computer: The `today` command can retrieve a random passage from the ESV
 * :computer: The `today` command can show a specific passage from the ESV
