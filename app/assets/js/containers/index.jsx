
var query = app.func.query('q'), clients = [], labels = [], candidates = [], labelOverrideNames = '',
    filters = {
      client: app.func.query('c', -1),
      status: app.storage.get('filters-status', 3),
      label: parseInt(app.func.query('l', '0'), 10),
      text: ''
    },
    isViewOnly = false;
if (query != '') {
  filters.text = query.replace(/\s/g,' ').replace(/　/g,' ');
  filters.text = filters.text.replace(/^\s+|\s+$/gm,'').toUpperCase();
}

$(document).ready(function () {
  $('#menu-containers').addClass('active');
  $('#container-detail pre').css({height: ($(window).height()-200)+'px'});
  labelOverrideNames = $('#override-label-id').val();
  isViewOnly = ($('#mode-view-only').val() == 'true');

  var search = $('#search-text').blur(_search);
  if (query != '') search.val(query);

  setStatusFilter(filters.status);
  $('#status-filter .dropdown-menu a').click(function(e) {
    setStatusFilter(parseInt($(this).attr('href').substring(1), 10));
    ReactDOM.render(<Table reload={true} />, document.getElementById('data'));
    app.func.stop(e);
  });
  $('.detail-refresh a').click(function (e) {
    _detail();
    app.func.stop(e);
  });
  $('#container-name').on('shown.bs.modal', function () {
    $('#container-name input').focus();
  });
  $('#container-name .act-rename').click(function (e) {
    var popup = $('#container-name'),
        name = popup.find('.container-name').val(),
        client = popup.find('.client-id').val(),
        newname = app.func.trim(popup.find('input.new-name').val());
    if (newname.length == 0) {
      popup.find('input.new-name').focus();
      app.func.stop(e);
      return;
    }
    popup.modal('hide');
    _rename(client, name, newname);
    app.func.stop(e);
  });
  $('#container-commit').on('shown.bs.modal', function () {
    $('#container-commit .tag').focus();
  });
  $('#container-commit .act-commit').click(function () {
    var popup = $('#container-commit'),
        name = popup.find('.container-name').val(),
        client = popup.find('.client-id').val(),
        repository = popup.find('.repository').val(),
        tag = popup.find('.tag').val(),
        message = popup.find('.message').val(),
        author = popup.find('.author').val();
    if (repository.length == 0) {
      popup.find('.repository').focus();
      return;
    }
    popup.modal('hide');
    _commit(client, name, repository, tag, message, author);
  });
});

$(window).keyup(function () {
  var search = $('#search-text');
  if (search.is(':focus')) {
    _search();
  }
});

function _setClientOption() {
  var options = $('.client-filters').hide(),
      caption = false,
      count = 0;
  options.find('ul.dropdown-menu').html('');
  $.map(clients, function (client) {
    if ((! caption) && (filters.client == client.key)) caption = client.value;
    count++;
  });
  caption && options.find('.caption').text(caption);
  if (count <= 1) {
    return;
  }
  var html = '<li><a href="#-1">All</a></li>';
  $.map(clients, function (client) {
    html += '<li><a href="#'+client.key+'">'+client.value+'</a></li>';
  });
  options.find('ul.dropdown-menu').html(html);
  options.fadeIn();

  $('#client-filter .dropdown-menu a').click(function(e) {
    var a = $(this),
        group = a.closest('.btn-group').removeClass('open');
    group.find('.caption').text(a.text()).blur();
    filters.client = a.attr('href').substring(1);
    ReactDOM.render(<Table reload={false} />, document.getElementById('data'));
    app.func.stop(e);
  });
}

function _setLabelFilter() {
  var options = $('.label-filters').hide(),
      caption = false,
      count = 0;
  options.find('ul.dropdown-menu').html('');
  $.map(labels, function (label) {
    if ((! caption) && (filters.label == app.func.hash(label.key+'->'+label.value))) {
      caption = label.value;
    }
    count++;
  });
  if ((! caption) && (filters.label == -1)) caption = 'Not Labeled';
  caption && options.find('.caption').text(caption);
  if (count <= 1) {
    return;
  }
  var html = '<li><a href="#0">All</a></li>',
      group = '';
  html += '<li><a href="#-1">Not Labeled</a></li>';
  $.map(labels, function (label) {
    if (group != label.key) {
      html += '<li class="dropdown-header">'+label.key+'</li>';
      group = label.key;
    }
    var value = label.value;
    if (value.length > 20) {
      value = value.substring(0, 20) + '..';
    }
    value = '&nbsp;&nbsp;&nbsp;&nbsp;'+value;
    html += '<li><a href="#'+app.func.hash(label.key+'->'+label.value)+'">'+value+'</a></li>';
  });
  options.find('ul.dropdown-menu').html(html);
  options.fadeIn();

  $('#label-filter .dropdown-menu a').click(function(e) {
    var a = $(this),
        group = a.closest('.btn-group').removeClass('open');
    group.find('.caption').text(a.text().trim()).blur();
    filters.label = a.attr('href').substring(1);

    if (window.history && window.history.pushState) {
      var url = (filters.label == 0) ? '/' : '/?l='+filters.label;
      history.pushState(null, null, url);
    }
    ReactDOM.render(<Table reload={false} />, document.getElementById('data'));
    app.func.stop(e);
  });
}

function setStatusFilter(value) {
  var a = $('#status-filter a[href="#'+value+'"]'),
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
  ReactDOM.render(<Table reload={false} />, document.getElementById('data'));
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
  }, error: function () {
    arg.err && alert(arg.err)
  }});
}

function _rename(client, id, newname) {
  var data = {name: newname};
  if (client) data.client = client;

  app.func.ajax({type: 'POST', url: '/api/container/rename/'+id, data: data, success: function (data) {
    if (data.indexOf('successfully') == -1) {
      alert(data);
      return;
    }
    ReactDOM.render(<Table reload={true} />, document.getElementById('data'));
  }});
}

function _commit(client, id, repository, tag, message, author) {
  var data = {repo: repository, tag: tag, msg: message, author: author};
  if (client) data.client = client;

  app.func.ajax({type: 'POST', url: '/api/container/commit/'+id, data: data, success: function (data) {
    var message = data.error ? data.error : 'committed successfully.';
    alert(message);
  }});
}

function _client(multiple, single) {
  return (clients.length > 1) ? multiple : single;
}

var TableRow = React.createClass({
  propTypes: {
    content: React.PropTypes.object.isRequired
  },
  inspect: function() {
    var tr = $(ReactDOM.findDOMNode(this)),
        id = tr.attr('data-container-id'),
        nm = '['+id.substring(0, 4)+'] '+tr.find('.dropdown a.dropdown-toggle').attr('data-container-name'),
        client = _client('?client='+tr.attr('data-client-id'), '');
    _detail({title: nm, url: '/api/container/inspect/'+id+client});
  },
  processes: function() {
    var tr = $(ReactDOM.findDOMNode(this)),
        name = tr.find('.dropdown a.dropdown-toggle').attr('data-container-name'),
        client = _client('?client='+tr.attr('data-client-id'), '');
    app.func.link('/container/top/'+name+client);
  },
  statlog: function() {
    var tr = $(ReactDOM.findDOMNode(this)),
        name = tr.find('.dropdown a.dropdown-toggle').attr('data-container-name'),
        client = _client('?client='+tr.attr('data-client-id'), '');
    app.func.link('/container/statlog/'+name+client);
  },
  changes: function() {
    var tr = $(ReactDOM.findDOMNode(this)),
        name = tr.find('.dropdown a.dropdown-toggle').attr('data-container-name'),
        client = _client('?client='+tr.attr('data-client-id'), '');
    app.func.link('/container/changes/'+name+client);
  },
  _action: function (sender, arg) {
    var tr = $(ReactDOM.findDOMNode(this)),
        name = tr.find('.dropdown a.dropdown-toggle').attr('data-container-name'),
        client = _client({client: tr.attr('data-client-id')}, '');
    var success = arg.success ? arg.success : function (data) {
      var message = data.error ? data.error : (arg.msg ? arg.msg : arg.action + 'ed') + ' successfully.';
      ReactDOM.render(<Table reload={true} />, document.getElementById('data'));
      alert(message);
    };
    if (arg.confirm && !window.confirm(arg.confirm + name)) {
      return;
    }
    app.func.ajax({type: 'POST', url: '/api/container/'+arg.action+'/'+name, data: client, success: success});
  },
  restart: function() {
    if (isViewOnly) return;
    this._action(this, {action: 'restart'});
  },
  start: function() {
    if (isViewOnly) return;
    this._action(this, {action: 'start'});
  },
  stop: function() {
    if (isViewOnly) return;
    this._action(this, {action: 'stop', msg: 'stopped'});
  },
  kill: function() {
    if (isViewOnly) return;
    this._action(this, {action: 'kill'});
  },
  rm: function() {
    if (isViewOnly) return;
    this._action(this, {action: 'rm', msg: 'removed', confirm: 'Are you sure to remove the container?\nID: '});
  },
  rename: function() {
    if (isViewOnly) return;
    var tr = $(ReactDOM.findDOMNode(this)),
        caption = tr.find('.dropdown a.dropdown-toggle').text(),
        name = tr.find('.dropdown a.dropdown-toggle').attr('data-container-name'),
        popup = $('#container-name');
    popup.find('.title').text(caption);
    popup.find('input').val(name);
    popup.find('.container-name').val(name);
    popup.find('.client-id').val(tr.attr('data-client-id'));
    popup.modal('show');
  },
  commit: function() {
    if (isViewOnly) return;
    var tr = $(ReactDOM.findDOMNode(this)),
        name = tr.find('.dropdown a.dropdown-toggle').text(),
        repo = tr.find('.image').text(),
        popup = $('#container-commit'),
        index = repo.indexOf(':');
    popup.find('.title').text(name);
    popup.find('.repository').val(((index == 0) ? repo : repo.substring(0, index)));
    popup.find('.tag').val(((index == 0) ? '' : repo.substring(index + 1)));
    popup.find('.container-name').val(name);
    popup.find('.client-id').val(tr.attr('data-client-id'));
    popup.modal('show');
  },
  image: function() {
    var tr = $(ReactDOM.findDOMNode(this)),
        image = tr.find('.image').text(),
        client = _client('&c='+tr.attr('data-client-id'), '');
    app.func.link('/images?q='+image+client);
  },
  render: function() {
    var props = this.props.content,
        container = props.container,
        key = container.id.substring(0, 20),
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
    if (labelOverrideNames && container.labels) {
      $.map(container.labels, function (value, key) {
        if (key == labelOverrideNames) {
          names = value + '..';
        }
      });
    }
    if (command.length > 13) {
      command = command.substring(0, 13) + '..';
    }
    status = status ? status.replace(/seconds*/, 'sec').replace(/minutes*/, 'min').replace(/About /, '') : '';
    if (isViewOnly) {
      return (
        <tr key={key + props.endpoint} data-client-id={props.client.id} data-container-id={key}>
          <td className="data-index">{container.id.substring(0, 4)}</td>
          <td className="data-name"><ul className="nav">
            <li className="dropdown">
              <a className="dropdown-toggle" data-toggle="dropdown" href="#" aria-expanded="true" data-container-name={name}>
                <span>{names.substring(0, names.length-2) + props.endpoint}</span>
              </a>
              <ul className="dropdown-menu">
                <li><a onClick={this.inspect}>inspect</a></li>
                <li><a onClick={this.processes}>processes (top)</a></li>
                <li><a onClick={this.statlog}>stats & logs</a></li>
                <li><a onClick={this.changes}>changes (diff)</a></li>
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
    } else {
      return (
        <tr key={key + props.endpoint} data-client-id={props.client.id} data-container-id={key}>
          <td className="data-index">{container.id.substring(0, 4)}</td>
          <td className="data-name"><ul className="nav">
            <li className="dropdown">
              <a className="dropdown-toggle" data-toggle="dropdown" href="#" aria-expanded="true" data-container-name={name}>
                <span>{names.substring(0, names.length-2) + props.endpoint}</span>
              </a>
              <ul className="dropdown-menu">
                <li><a onClick={this.inspect}>inspect</a></li>
                <li><a onClick={this.processes}>processes (top)</a></li>
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
  }
});

var Table = React.createClass({
  propTypes: {
    reload: React.PropTypes.bool.isRequired
  },
  getInitialState: function() {
    return {data: []};
  },
  load: function(sender) {
    clients = [];
    labels = [];

    app.func.ajax({type: 'GET', url: '/api/containers', data: {status: filters.status}, success: function (data) {
      if (data.error) {
        alert(data.error);
        return;
      }
      candidates = data;

      // make filters
      var temp = {clients: {}, labels: {}},
          conf = $('#filter-label-ids').val();
      $.map(candidates, function (candidate) {
        temp.clients[''+candidate.client.id] = candidate.client.endpoint.replace(/^.*:\/\//, '').replace(/:.*$/, '');

        $.map(candidate.containers, function (container) {
          if (! container.labels) return;
          $.map(container.labels, function (value, key) {
            if ((conf != 'all') && (conf.indexOf(key) == -1)) return;
            if (! temp.labels[key]) {
              temp.labels[key] = {};
            }
            temp.labels[key][value] = true;
          });
        });
      });
      $.map(temp.clients, function (value, key) {
        clients.push({key: key, value: value});
      });
      $.map(temp.labels, function (nest, key) {
        $.map(nest, function (_, value) {
          labels.push({key: key, value: value});
        });
      });
      clients.sort(function (a, b) {
        if (a.value < b.value) return -1;
        if (a.value > b.value) return 1;
        return 0;
      });
      labels.sort(function (a, b) {
        if (a.key < b.key) return -1;
        if (a.key > b.key) return 1;
        if (a.value < b.value) return -1;
        if (a.value > b.value) return 1;
        return 0;
      });
      _setClientOption();
      _setLabelFilter();

      // reflow
      sender.setState({data: sender.filter()});
    }});
  },
  filter: function() {
    var data = [];
    $.map(candidates, function (candidate) {
      if ((filters.client == -1) || (candidate.client.id == filters.client)) {
        $.map(candidate.containers, function (container) {
          if ((filters.label == 0) || ((filters.label == -1) && (! container.labels))) {
            data.push({client: candidate.client, container: container});
          } else {
            if (! container.labels) return;
            var match = false;
            $.map(container.labels, function (value, key) {
              match |= (filters.label == app.func.hash(key+'->'+value));
            });
            if (match) data.push({client: candidate.client, container: container});
          }
        });
      }
    });
    data.sort(function (a, b) {
      if (a.container.names.join(',') < b.container.names.join(','))
        return -1;
      if (a.container.names.join(',') > b.container.names.join(','))
        return 1;
      return 0;
    });
    $('#count').text(data.length + ' container' + ((data.length > 1) ? 's' : ''));
    return data;
  },
  componentDidMount: function() {
    this.load(this);
  },
  componentWillReceiveProps: function(props) {
    if (props.reload) {
      this.load(this);
      return;
    }
    this.setState({data: this.filter()});
  },
  render: function() {
    var multiple = (candidates.length > 1) && (filters.client == -1);

    var rows = this.state.data.map(function(record, index) {
      var client = record.client,
          container = record.container,
          key = record.container.id.substring(0, 10);

      if (filters.text != '') {
        var match = true;
        $.map(filters.text.split(' '), function (word) {
          var innerMatch = (container.id.substring(0, 10).toUpperCase().indexOf(word) > -1);
          $.map(container.names, function (name) {
            innerMatch |= (name.toUpperCase().indexOf(word) > -1);
          });
          innerMatch |= (container.image.toUpperCase().indexOf(word) > -1);
          innerMatch |= (container.command && (container.command.toUpperCase().indexOf(word) > -1));
          innerMatch |= (container.status && (container.status.toUpperCase().indexOf(word) > -1));
          if (container.ports) {
            $.map(container.ports, function (port) {
              innerMatch |= (port.Type.toUpperCase().indexOf(word) > -1);
              innerMatch |= (port.IP.indexOf(word) > -1);
              innerMatch |= ((''+port.PrivatePort).indexOf(word) > -1);
              innerMatch |= ((''+port.PublicPort).indexOf(word) > -1);
            });
          }
          match &= innerMatch;
        });
        if (! match) return;
      }
      return <TableRow key={key+'@'+client.id} index={key+'@'+index} content={{
          endpoint: multiple ? ' @'+client.endpoint.replace(/^.*:\/\//, '').replace(/:.*$/, ''): '',
          client: client, container: container
      }} />
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

ReactDOM.render(<Table reload={false} />, document.getElementById('data'));
