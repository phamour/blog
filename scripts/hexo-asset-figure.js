'use strict';

const { resolve } = require('url');
const { encodeURL, htmlTag } = require('hexo-util');
const img = require('hexo/lib/plugins/tag/img')(hexo);

const PostAsset = hexo.model('PostAsset');
const rMeta = /title\s*=\s*["']?([^"']+)?["']?/;

/**
 * Asset figure tag
 *
 * Syntax:
 *   {% asset_figure [class names] slug [width] [height] [title text [alt text]]%}
 */
hexo.extend.tag.unregister('asset_figure');
hexo.extend.tag.register('asset_figure', function (args) {
  // Find image URL
  const len = args.length;
  for (let i = 0; i < len; i++) {
    const asset = PostAsset.findOne({post: this._id, slug: args[i]});
    if (asset) {
      args[i] = encodeURL(resolve('/', asset.path));
    }
  }

  // Build image tag
  const imgTag = img(args);
  
  // Find title from image tag
  let title;
  if (imgTag) {
    const match = rMeta.exec(imgTag);
    if (match != null) {
      title = match[1];
    }
  }

  // Build caption tag and the final figure tag
  const captionTag = htmlTag('figcaption', {}, title)
  const figureTag = htmlTag('figure', {}, imgTag + captionTag, false);
  return figureTag;
});
