$(function() {
  moment.locale('sv');

  var feedTemplate = [
    '<div class="panel panel-default">',
    '<div class="panel-heading">',
    '<h3 class="panel-title">{{> title}}</h3>',
    '</div>',
    '<div class="panel-body">',
    '<ul class="list-group">',
    '{{for items}}',
    '<li class="list-group-item"><strong>{{>updated}}</strong> <a href="{{:link}}">{{>title}}</a> <small>{{:author}}</small>',
    '<span class="pull-right"><img src="{{:thumbnail}}" width="30" height="30"></span>',
    '</li>',
    '{{/for}}',
    '</ul>',
    '</div>',
    '</div>'
  ].join('');

  $.templates("feedTemplate", feedTemplate);

  $.getJSON("/feed/list", {}, function(feeds) {
    feeds.sort(function(a, b) {
      return a.Name.localeCompare(b.Name);
    });
    $.each(feeds, function(i, feed) {
      $('.feedlist').append('<div class="feed-' + feed.Name + '"/>');
    });

    $.each(feeds, function(i, feed) {
      $.get("/feed/" + feed.Name, function (data) {
        var items = [];
        $(data).find("entry").each(function () { // or "item" or whatever suits your feed
          var el = $(this);
          var item = {
            title: el.find("title").text(),
            author: el.find("author name").text(),
            thumbnail: el.find("media\\:thumbnail,thumbnail").attr('url'),
            updated: moment(el.find("updated").text()).format("LLL"),
            link: el.find("link").attr('href'),
          };
          items.push(item);
        });
        var html = $.render.feedTemplate({
          title: feed.Name,
          items: items.slice(0, 5)
        });
        $('.feedlist .feed-' + feed.Name).html(html);
      });
    });
  });
});
