
var table, query = app.func.query('q'),
    filters = {status: 0, text: ''};
if (query != '') {
  filters.text = query.replace(/\s/g,' ').replace(/　/g,' ');
  filters.text = filters.text.replace(/^\s+|\s+$/gm,'').toUpperCase();
}

$(document).ready(function () {
  $('#menu-containers').addClass('active');
  $('#container-detail pre').css({height: ($(window).height()-200)+'px'})

  var search = $('#search-text').blur(_search);
  if (query != '') search.val(query);

  $('#conditions .dropdown-menu a').click(function() {
    var a = $(this), group = a.closest('.btn-group').removeClass('open');
    filters[group.attr('data-filter-key')] = parseInt(a.attr('href').substring(1), 10);
    group.find('.caption').text(a.text()).blur();
    table.setProps();
    return false;
  });
  $('.detail-refresh a').click(function (e) {
    _detail();
    return false;
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
  var popup = $('#container-detail'),
      details = popup.find('.details').hide();
  app.func.ajax({type: 'GET', url: arg.url, data: arg.conditions, success: function (data) {
    popup.find('.detail-title').text(arg.title);
    details.text(arg.format(data)).fadeIn();
    popup.modal('show');
    last = arg;
  }, error: function (xhr, status, err) {
    arg.err && alert(arg.err)
  }});
}

var TableRow = React.createClass({
  inspect: function() {
    var tr = $(this.getDOMNode()),
        id = tr.attr('data-container-id'),
        nm = '['+id.substring(0, 4)+'] '+tr.find('.dropdown a.dropdown-toggle').attr('data-container-name');
    _detail({title: nm, url: '/api/container/inspect/'+id});
    return false;
  },
  processes: function() {
    var name = $(this.getDOMNode()).find('.dropdown a.dropdown-toggle').attr('data-container-name');
    location.href = '/container/top/'+name;
    return false;
  },
  statlog: function() {
    var name = $(this.getDOMNode()).find('.dropdown a.dropdown-toggle').attr('data-container-name');
    location.href = '/container/statlog/'+name;
    return false;
  },
  changes: function() {
    var name = $(this.getDOMNode()).find('.dropdown a.dropdown-toggle').attr('data-container-name');
    location.href = '/container/changes/'+name;
    return false;
  },
  start: function() {
    var name = $(this.getDOMNode()).find('.dropdown a.dropdown-toggle').attr('data-container-name');
    app.func.ajax({type: 'POST', url: '/api/container/start/'+name, success: function (data) {
      if (data.error) {
        alert(data.error);
        return;
      }
      table.setProps();
    }});
    return false;
  },
  stop: function() {
    var name = $(this.getDOMNode()).find('.dropdown a.dropdown-toggle').attr('data-container-name');
    app.func.ajax({type: 'POST', url: '/api/container/stop/'+name, success: function (data) {
      if (data.error) {
        alert(data.error);
        return;
      }
      table.setProps();
    }});
    return false;
  },
  restart: function() {
    var name = $(this.getDOMNode()).find('.dropdown a.dropdown-toggle').attr('data-container-name');
    app.func.ajax({type: 'POST', url: '/api/container/restart/'+name, success: function (data) {
      if (data.error) {
        alert(data.error);
        return;
      }
      table.setProps();
    }});
    return false;
  },
  rm: function() {
    var name = $(this.getDOMNode()).find('.dropdown a.dropdown-toggle').attr('data-container-name');
    if (!window.confirm('Are you sure to remove a container: '+name)) {
      return;
    }
    app.func.ajax({type: 'POST', url: '/api/container/rm/'+name, success: function (data) {
      if (data != 'removed successfully.') {
        alert(data);
        return;
      }
      table.setProps();
    }});
    return false;
  },
  image: function() {
    var name = $(this.getDOMNode()).find('.image').text();
    location.href = '/images?q='+name;
    return false;
  },
  render: function() {
    var container = this.props.content,
        names = '', name = '', ports = '',
        command = container.command,
        status = container.status;
    if (container.ports) {
      $.map(container.ports, function (port) {
        ports += (port.IP ? port.IP+':' : '')+(port.PublicPort ? port.PublicPort+'->' : '')+port.PrivatePort+'/'+port.Type+',';
      });
    }
    if (container.names) {
      $.map(container.names, function (n) {
        names += n.replace('/', '') + ',';
        name = n.replace('/', '');
      });
    }
    if (command.length > 13) {
      command = command.substring(0, 13) + '..';
    }
    status = status ? status.replace(/seconds*/, 'sec').replace(/minutes*/, 'min').replace(/About /, '') : '';
    return (
        <tr key={this.props.index} data-container-id={container.id.substring(0, 20)}>
          <td className="data-index">{container.id.substring(0, 4)}</td>
          <td className="data-name"><ul className="nav">
            <li className="dropdown">
              <a className="dropdown-toggle" data-toggle="dropdown" href="#" aria-expanded="true" data-container-name={name}>
                <span>{names.substring(0, names.length-1)}</span>
              </a>
              <ul className="dropdown-menu">
                <li><a onClick={this.inspect}>inspect</a></li>
                <li><a onClick={this.processes}>processes</a></li>
                <li><a onClick={this.statlog}>stats & logs</a></li>
                <li><a onClick={this.changes}>changes (diff)</a></li>
                <li className="divider"></li>
                <li><a onClick={this.start}>start</a></li>
                <li><a onClick={this.stop}>stop</a></li>
                <li><a onClick={this.restart}>restart</a></li>
                <li><a onClick={this.rm}>rm</a></li>
                <li className="divider"></li>
                <li><a onClick={this.image}>image</a></li>
              </ul>
            </li>
          </ul></td>
          <td className="data-name">{command}</td>
          <td className="data-name">{ports.substring(0, ports.length-1)}</td>
          <td className="data-name image">{container.image}</td>
          <td className="data-name">{status}</td>
        </tr>
    );
  }
});

var Table = React.createClass({
  getInitialState: function() {
    return {data: []};
  },
  load: function(sender) {
    var conditions = {
      status: filters.status,
      q: filters.text
    };
    app.func.ajax({type: 'GET', url: '/api/containers', data: conditions, success: function (data) {
      $('#count').text(data.length + ' container' + ((data.length > 1) ? 's' : ''));
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
          <TableRow key={record.id.substring(0, 10)} index={index} content={record} />
      );
    });
    return (
        <table className="table table-striped table-hover">
          <thead>
            <tr>
              <th>ID</th>
              <th>Names</th>
              <th>Command</th>
              <th>Ports</th>
              <th>Repository & Tags</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>{rows}</tbody>
        </table>
    );
  }
});

table = React.render(<Table />, document.getElementById('data'));
