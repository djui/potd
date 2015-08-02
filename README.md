# potd - Paper of the Day

Scrape `@onepaperperday` and `@onecspaperaday` on Twitter and print the latest
paper tweet.


# Usage

    $ ./potd
    By default Windows 10 tracks/shares yr clicks, purchases, places, typing
    http://t.co/aqzhgZDAq3
    http://www.polygon.com/2015/7/31/9075531/windows-10-privacy-how-to


# Build

    $ go build

A `CREDENTIALS` file next to the executable with Twitter API credentials is required.

    $ echo -e "$CONSUMER_KEY\n$CONSUMER_SECRET\n\n\n" > CREDENTIALS

