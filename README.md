# CardDeck

CardDeck is a simple web server app for choosing random cards from decks of card images. It uses `embed` to make a
single binary with all the card images, so it's trivial to distribute.

The original intention is to provide an online interface for Nathan Rockwood's
[GM's Apprentice](https://www.drivethrurpg.com/product/125685/The-GameMasters-Apprentice-Base-Deck) decks. Any collection of images will do, however.

Once the server is running, open it in a browser at http://localhost:8888 . You can select any installed deck to draw a random card.

To add your own decks at build time, add a directory under `decks` with a short descriptive name. This name will be a button on
the web page. In that directory place the individual jpeg images for your card deck. They will be automatically included when you `go build`.

If you only have the binary or don't want to embed the decks, you can load decks at runtime instead. In the same directory
where you have the `carddeck` binary, create a directory `decks`. Put your cards there with each in its own directory,
like the sample Poker and Tarot decks.

If local decks exist they will be used instead of embedded decks. To use embedded decks anyway, use the command-line option `carddeck --embedded`.

The server runs on port 8888 by default, but you can use the command line argument `-port=XXXX` to change this.

To build and run the server as a Docker container:
* docker build -t carddeck
* docker run --rm -p 8888:8888 carddeck

This is tested with Docker Desktop on OSX. Obviously you can get a lot more complex with Docker, like mounting the cards
separately in the image or a bind volume, but if you know you can do that you probably know how to do it too. In any
event, running an application in Docker that's already optimized to be a single binary with no state is kind of gratuitous,
so it's not very important that it be perfect.

TODO:
* Store the cards a user has seen in a session, so they don't get repeats without asking to reshuffle the decks
* Add documentation links like an image of a sample card from the GMA manual
* Allow loading decks from zip files instead of open directories

Credits:
* Tarot deck from [Luciella Elisabeth Scarlett](https://luciellaes.itch.io/rider-waite-smith-tarot-cards-cc0) (Public Domain and CC0 license)
* Poker deck from [Byron Knoll](http://byronknoll.blogspot.com/2011/03/vector-playing-cards.html) (Public Domain)
