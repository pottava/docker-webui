
var table, query = app.func.query('q'),
    filters = {status: app.storage.get('filters-status', 3), text: ''};
if (query != '') {
  filters.text = query.replace(/\s/g,' ').replace(/　/g,' ');
  filters.text = filters.text.replace(/^\s+|\s+$/gm,'').toUpperCase();
}

$(document).ready(function () {
  $('#menu-containers').addClass('active');
  $('#container-detail pre').css({height: ($(window).height()-200)+'px'})

  var search = $('#search-text').blur(_search);
  if (query != '') search.val(query);

  setStatusFilter(filters.status);
  $('#conditions .dropdown-menu a').click(function() {
    setStatusFilter(parseInt($(this).attr('href').substring(1), 10));
    table.setProps();
    return false;
  });
  $('.detail-refresh a').click(function (e) {
    _detail();
    return false;
  });
  $('#container-name').on('shown.bs.modal', function (e) {
    $('#container-name input').focus();
  });
  $('#container-name .act-rename').click(function (e) {
    var popup = $('#container-name'),
        name = popup.find('.title').text(),
        newname = app.func.trim(popup.find('input').val());
    if (newname.length == 0) {
      popup.find('input').focus();
      return;
    }
    popup.modal('hide');
    _rename(name, newname);
  });
  $('#container-commit').on('shown.bs.modal', function (e) {
    $('#container-commit .tag').focus();
  });
  $('#container-commit .act-commit').click(function (e) {
    var popup = $('#container-commit'),
        name = popup.find('.title').text(),
        repository = popup.find('.repository').val(),
        tag = popup.find('.tag').val(),
        message = popup.find('.message').val(),
        author = popup.find('.author').val();
    if (repository.length == 0) {
      popup.find('.repository').focus();
      return;
    }
    popup.modal('hide');
    _commit(name, repository, tag, message, author);
  });
});

$(window).keyup(function (e) {
  var search = $('#search-text');
  if (search.is(':focus') && (e.which == 13)) {
    _search();
  }
});

function setStatusFilter(value) {
  var a = $('#conditions a[href="#'+value+'"]'),
      group = a.closest('.btn-group').removeClass('open');
  group.find('.caption').text(a.text()).blur();
  app.storage.set('filters-status', value);
  filters.status = value;
}

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

function _rename(id, newname) {
  var data = {name: newname};
  app.func.ajax({type: 'POST', url: '/api/container/rename/'+id, data: data, success: function (data) {
    if (data.indexOf('successfully') == -1) {
      alert(data);
      return;
    }
    table.setProps();
  }});
}

function _commit(id, repository, tag, message, author) {
  var data = {repo: repository, tag: tag, msg: message, author: author};
  app.func.ajax({type: 'POST', url: '/api/container/commit/'+id, data: data, success: function (data) {
    var message = data.error ? data.error : 'committed successfully.';
    alert(message);
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
  _action: function (sender, arg) {
    var name = $(sender.getDOMNode()).find('.dropdown a.dropdown-toggle').attr('data-container-name');
    var success = arg.success ? arg.success : function (data) {
      var message = data.error ? data.error : (arg.msg ? arg.msg : arg.action + 'ed') + ' successfully.';
      table.setProps();
      alert(message);
    };
    if (arg.confirm && !window.confirm(arg.confirm + name)) {
      return;
    }
    app.func.ajax({type: 'POST', url: '/api/container/'+arg.action+'/'+name, success: success});
  },
  restart: function() {
    this._action(this, {action: 'restart'});
    return false;
  },
  start: function() {
    this._action(this, {action: 'start'});
    return false;
  },
  stop: function() {
    this._action(this, {action: 'stop', msg: 'stopped'});
    return false;
  },
  kill: function() {
    this._action(this, {action: 'kill'});
    return false;
  },
  rm: function() {
    this._action(this, {action: 'rm', msg: 'removed', confirm: 'Are you sure to remove the container?\nID: '});
    return false;
  },
  exec: function() {
    return false;
  },
  rename: function() {
    var name = $(this.getDOMNode()).find('.dropdown a.dropdown-toggle').attr('data-container-name');
    var popup = $('#container-name');
    popup.find('.title').text(name);
    popup.find('input').val(name);
    popup.modal('show');
    return false;
  },
  commit: function() {
    var tr = $(this.getDOMNode()),
        name = tr.find('.dropdown a.dropdown-toggle').attr('data-container-name'),
        repo = tr.find('.image').text(),
        popup = $('#container-commit');
    popup.find('.title').text(name);
    popup.find('.repository').val(repo.substring(0, repo.indexOf(':')));
    popup.find('.tag').val(repo.substring(repo.indexOf(':') + 1));
    popup.modal('show');
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
        names += n.replace('/', '') + ', ';
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
                <span>{names.substring(0, names.length-2)}</span>
              </a>
              <ul className="dropdown-menu">
                <li><a onClick={this.inspect}>inspect</a></li>
                <li><a onClick={this.processes}>processes</a></li>
                <li><a onClick={this.statlog}>stats & logs</a></li>
                <li><a onClick={this.changes}>changes (diff)</a></li>
                <li><a onClick={this.rename}>rename</a></li>
                <li className="divider"></li>
                <li><a onClick={this.restart}>restart</a></li>
                <li><a onClick={this.start}>start</a></li>
                <li><a onClick={this.stop}>stop</a></li>
                <li><a onClick={this.kill}>kill</a></li>
                <li><a onClick={this.rm}>rm</a></li>
                <li className="divider"></li>
                <li><a onClick={this.commit}>commit</a></li>
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
