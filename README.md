# CardDeck

CardDeck is a simple web server app for choosing random cards from decks of card images. It uses `embed` to make a
single binary with all the card images, so it's trivial to distribute.

The original intention is to provide an online interface for Nathan Rockwood's
[GM's Apprentice](https://www.drivethrurpg.com/product/125685/The-GameMasters-Apprentice-Base-Deck) decks. Any collection of images will do, however.

Once the server is running, open it in a browser at http://localhost:8888 . You can select any installed deck to draw a random card.

To add your own decks at build time, create a directory under `static` with a short descriptive name. This name will be a button on
the web page. In that directory place the individual jpeg images for your card deck. They will be automatically included when you `go build`.

To instead have your own cards at runtime, create the same directory `static` in the same directory as the binary,
put your card directories in it in the same manner, and run `carddeck live`. That instructs the program to load
cards from the filesystem instead of using any embedded decks.

The server runs on port 8888 by default but uses environment variable PORT to configure this.

To build and run the server as a Docker container:
* docker build -t carddeck
* docker run --rm -p 8888:8888 carddeck

This is tested with Docker Desktop on OSX. Obviously you can get a lot more complex with Docker like mounting the cards
separately in the image or a bind volume, but if you know you can do that you probably know how to do it too. In any
event, running an application in Docker that's already optimized to be a single binary with no state is kind of gratuitous,
so it's not very important that it be perfect.

TODO:
* Store the cards a user has seen in a session, so they don't get repeats without asking to reshuffle the decks
* Add documentation links like an image of a sample card from the GMA manual

Fun but pointless changes:
* Rewrite using a framework like Echo. Any framework is overkill for this because it just serves one template, 
one static CSS file, and a FS full of images, but it would be a useful learning project.

Credits:
* Tarot deck from [Luciella Elisabeth Scarlett](https://luciellaes.itch.io/rider-waite-smith-tarot-cards-cc0) (Public Domain and CC0 license)
* Poker deck from [Byron Knoll](http://byronknoll.blogspot.com/2011/03/vector-playing-cards.html) (Public Domain)
* I didn't use them but this deck generator looks nice: https://www.me.uk/cards/
