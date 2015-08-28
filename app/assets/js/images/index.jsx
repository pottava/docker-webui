
var table, query = app.func.query('q'),
    filters = {text: ''}, reload = false;
if (query != '') {
  filters.text = query.replace(/\s/g,' ').replace(/　/g,' ');
  filters.text = filters.text.replace(/^\s+|\s+$/gm,'').toUpperCase();
}

$(document).ready(function () {
  $('#menu-images').addClass('active');
  $('#image-detail pre').css({height: ($(window).height()-200)+'px'})

  var search = $('#search-text').blur(_search);
  if (query != '') search.val(query);

  $('.detail-refresh a').click(function (e) {
    _detail();
    return false;
  });
  $('#image-detail').on('hide.bs.modal', function (e) {
    if (reload) table.setProps();
  });
  $('#image-pull').on('shown.bs.modal', function (e) {
    $('#image-name').val('').focus();
  });
  $('#image-pull .act-pull').click(function (e) {
    var name = app.func.trim($('#image-name').val());
    if (name.length == 0) {
      $('#image-name').focus();
      return;
    }
    $('#image-pull').modal('hide');
    _pull(name);
  });
  $('#image-run .act-run').click(function (e) {
    $('#image-run').modal('hide');
  });
});

$(window).keyup(function (e) {
  var search = $('#search-text');
  if (search.is(':focus') && (e.which == 13)) {
    _search();
  }
});

function _search() {
  var candidate = $('#search-text').val().replace(/\s/g,' ').replace(/　/g,' ');
  candidate = candidate.replace(/^\s+|\s+$/gm,'').toUpperCase();
  if (filters.text == candidate) return;
  filters.text = candidate;
  table.setProps();
}

var last = {};

function _detail(arg) {
  arg = arg ? arg : last;
  arg.format = arg.format ? arg.format : function (data) {
    return JSON.stringify(data, true, ' ');
  };
  var popup = $('#image-detail'),
      details = popup.find('.details');
  popup.find('.detail-title').text(arg.title);
  popup.find('.detail-refresh').hide();
  if (arg.message) {
    details.text(arg.message);
  } else {
    details.hide();
  }
  app.func.ajax({type: 'GET', url: arg.url, data: arg.conditions, success: function (data) {
    var formatted = arg.format(data)
    if (formatted.indexOf('Error:') == -1) {
      popup.find('.detail-refresh').show();
      details.text(formatted).fadeIn();
    } else {
      details.text(data).fadeIn();
    }
    popup.modal('show');
    last = arg;
  }, error: function (xhr, status, err) {
    arg.err && alert(arg.err)
  }});
}

function _pull(name) {
  reload = true;
  $('#image-detail').modal('show');
  _detail({title: name, message: 'Now executing..\n\ndocker pull '+name, url: '/api/image/pull/'+name});
}

var TableRow = React.createClass({
  inspect: function() {
    var tr = $(this.getDOMNode()),
        id = tr.attr('data-image-id'),
        nm = tr.find('.dropdown a.dropdown-toggle').text();
    _detail({title: nm, url: '/api/image/inspect/'+id});
    return false;
  },
  history: function() {
    var name = $(this.getDOMNode()).find('.dropdown a.dropdown-toggle').text();
    location.href = '/image/history/'+name;
    return false;
  },
  containers: function() {
    var name = $(this.getDOMNode()).find('.dropdown a.dropdown-toggle').text();
    location.href = '/?q='+name;
    return false;
  },
  pull: function() {
    _pull($(this.getDOMNode()).find('.dropdown a.dropdown-toggle').text())
    return false;
  },
  rmi: function() {
    var name = $(this.getDOMNode()).find('.dropdown a.dropdown-toggle').text();
    if (!window.confirm('Are you sure to remove image: '+name)) {
      return;
    }
    app.func.ajax({type: 'POST', url: '/api/image/rmi/'+name, success: function (data) {
      if (data != 'removed successfully.') {
        alert(data);
        return;
      }
      table.setProps();
    }});
    return false;
  },
  run: function() {
    var name = $(this.getDOMNode()).find('.dropdown a.dropdown-toggle').text(),
        popup = $('#image-run');
    $('#run-scripts').val('docker run ' + name);
    popup.find('.detail-title').text('Run from ' + name);
    popup.modal('show');
    return false;
  },
  render: function() {
    var image = this.props.content;
    return (
        <tr key={this.props.index} data-image-id={image.id.substring(0, 20)}>
          <td className="data-index">{image.id.substring(0, 10)}</td>
          <td className="data-name"><ul className="nav">
            <li className="dropdown">
              <a className="dropdown-toggle" data-toggle="dropdown" href="#" aria-expanded="true">{image.repoTags}</a>
              <ul className="dropdown-menu">
                <li><a onClick={this.inspect}>inspect</a></li>
                <li><a onClick={this.history}>history</a></li>
                <li className="divider"></li>
                <li><a onClick={this.containers}>containers</a></li>
                <li className="divider"></li>
                <li><a onClick={this.pull}>pull again</a></li>
                <li><a onClick={this.rmi}>rmi</a></li>
              </ul>
            </li>
          </ul></td>
          <td className="data-name">{app.func.byteFormat(image.virtualSize)}</td>
          <td className="data-name">{app.func.relativeTime(new Date(image.created * 1000))}</td>
        </tr>
    );
  }
});

var Table = React.createClass({
  getInitialState: function() {
    return {data: []};
  },
  load: function(sender) {
    var conditions = {q: filters.text};
    app.func.ajax({type: 'GET', url: '/api/images', data: conditions, success: function (data) {
      $('#count').text(data.length + ' image' + ((data.length > 1) ? 's' : ''));
      sender.setState({data: data});
    }});
  },
  componentDidMount: function() {
    this.load(this);
  },
  componentWillReceiveProps: function() {
    this.load(this);
  },
  render: function() {
    var rows = this.state.data.map(function(record, index) {
      return (
          <TableRow key={record.name} index={index} content={record} />
      );
    });
    return (
        <table className="table table-striped table-hover">
          <thead>
            <tr>
              <th>ID</th>
              <th>Repository & Tags</th>
              <th>VirtualSize</th>
              <th>Created</th>
            </tr>
          </thead>
          <tbody>{rows}</tbody>
        </table>
    );
  }
});

table = React.render(<Table />, document.getElementById('data'));
