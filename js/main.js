$(function() {
  moment.locale('sv');

  var feedTemplate = [
    '<div class="col-md-12">',
    '<h2>{{> title}}</h2>',
    '<ul class="list-unstyled">',
    '{{for items}}',
    '<li><a href="{{:link}}">{{>title}}</a> <small>{{>updated}}</small></li>',
    '{{/for}}',
    '</ul>',
    '</div>'
  ].join('');

  $.templates("feedTemplate", feedTemplate);

  var feeds = [
    "haproxy",
    "crmsh",
    "hawk",
    "resource-agents",
    "pacemaker",
    "fence-agents"
  ];

  $.each(feeds, function(i, name) {
    $.get("/feed/" + name, function (data) {
      var items = [];
      $(data).find("entry").each(function () { // or "item" or whatever suits your feed
        var el = $(this);
        var item = {
          title: el.find("title").text(),
          author: el.find("author").text(),
          updated: moment(el.find("updated").text()).format("LLL"),
          link: el.find("link").text(),
        };
        items.push(item);
      });
      var html = $.render.feedTemplate({
        title: name,
        items: items.slice(0, 5)
      });
      $('.feedlist').append(html);
    });
  });
});
