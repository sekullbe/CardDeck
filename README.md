# CardDeck

CardDeck is a simple web server app for choosing random cards from decks of card images. It uses `embed` to make a
single binary with all the card images, so it's trivial to distribute.

The original intention is to provide an online interface for Nathan Rockwood's
[GM's Apprentice](https://www.drivethrurpg.com/product/125685/The-GameMasters-Apprentice-Base-Deck) decks. Any colletion of images will do, however.

To add your own decks, create a directory under `static` with a short descriptive name. This name will be a button on
the web page. In that directory place the individual jpeg images for your card deck. They will be automatically included when you `go build`.


TODO:
* Add configurability to use local or zip FS instead of embed.FS
* Spiff up the interface
