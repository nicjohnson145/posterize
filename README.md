# posterize

Split a large PNG into smaller PNGs for printing.

## Why?

I make my campaign maps in [DungeonDraft](https://dungeondraft.net/). I want to make my maps all at
once and then split them into (in my case) multiple 8.5x11 sheets to make larger battle maps. I
couldnt find a solution that could fit my needs. Hence, `posterize`.


## Basic Usage

* Create your map in DungoenDraft (or whatever map making software you prefer).
* Export it as a PNG, making note of the PPI - Pixels Per Inch of the image.
    * If you're slicing a pre-made map that you dont know the PPI, but it does have grid, you can
      get a rough idea of the PPI by using the measure tool in something like Gimp and measuring one
      grid square
* `posterize <path_to_your_image> --pixels-per-inch <your PPI value>`
* Posterize should output your image split across multiple pages. Cut em out (there's usually a
  slight border since most printers cant print all the way to the edge). And either tape them together
  or use something like these [Acrylic Sign Holders](https://www.amazon.com/gp/product/B07F9SBW6H)
* Done!

## Examples

See the `examples` directory for example results, but [guildpact-archives.png](examples/guildpact-archives.png) is
the original image (2000px x 1600px @ 100ppi), and [img-0.png](examples/img-0.png), [img-1.png](examples/img-1.png), 
[img-2.png](examples/img-2.png), and [img-3.png](examples/img-3.png) are the resulting split images.
