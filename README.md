# CardDeck

CardDeck is a simple web server app for choosing random cards from decks of card images. It uses `embed` to make a
single binary with all the card images, so it's trivial to distribute.

The original intention is to provide an online interface for Nathan Rockwood's
[GM's Apprentice](https://www.drivethrurpg.com/product/125685/The-GameMasters-Apprentice-Base-Deck) decks. Any collection of images will do, however.

To add your own decks at build time, create a directory under `static` with a short descriptive name. This name will be a button on
the web page. In that directory place the individual jpeg images for your card deck. They will be automatically included when you `go build`.

To instead have your own cards at runtime, create the same directory `static` in the same directory as the binary,
put your card directories in it in the same manner, and run `carddeck live`. That instructs the program to load
cards from the filesystem instead of using any embedded decks.

TODO:
* Spiff up the interface
* Store the cards a user has seen in a session, so they don't get repeats without asking to reshuffle the decks


