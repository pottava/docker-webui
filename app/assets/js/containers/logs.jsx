var filters = {
      client: app.func.query('c', -1),
      label: parseInt(app.func.query('l', '0'), 10)
    },
    allcontainers = [],
    refreshWindow = app.storage.get('refresh-window-logs', 1),
    monitoringCount = app.storage.get('monitoring-count-logs', 20),
    stream = [];

$(document).ready(function () {
  $('#menu-logs').addClass('active');

  setRefreshWindow(refreshWindow);
  $('#refresh-window a').click(function(e) {
    setRefreshWindow(parseInt($(this).attr('href').substring(1), 10));
    app.func.stop(e);
  });
  setMonitoringCount(monitoringCount);
  $('#monitoring-count a').click(function(e) {
    setMonitoringCount(parseInt($(this).attr('href').substring(1), 10));
    app.func.stop(e);
  });

  // make filters
  app.func.ajax({url: '/api/containers', data: {status: 3}, success: function (data) {
    allcontainers = data;
    var clients = [],
        labels = [];

    var temp = {clients: {}, labels: {}},
        conf = $('#filter-label-ids').val();
    $.map(allcontainers, function (candidate) {
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
    _setClientOption(clients);
    _setLabelFilter(labels);
    refresh();
  }});
});

function setRefreshWindow(value) {
  var a = $('#refresh-window a[href="#'+value+'"]'),
      group = a.closest('.btn-group').removeClass('open');
  refreshWindow = value;
  app.storage.set('refresh-window-logs', value);
  group.find('.caption').text('refresh / '+a.text()).blur();
}
function setMonitoringCount(value) {
  var a = $('#monitoring-count a[href="#'+value+'"]'),
      group = a.closest('.btn-group').removeClass('open');
  monitoringCount = value;
  app.storage.set('monitoring-count-logs', value);
  group.find('.caption').text(a.text()).blur();
}

function refresh() {
  if (refreshWindow && (refreshWindow > 0)) {
    ReactDOM.render(<Table reload={true} />, document.getElementById('data'));
  }
  setTimeout(refresh, Math.max(1, refreshWindow) * 1000);
}

function _setClientOption(clients) {
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

function _setLabelFilter(labels) {
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
    ReactDOM.render(<Table reload={false} />, document.getElementById('data'));

    if (window.history && window.history.pushState) {
      var url = (filters.label == 0) ? '/logs' : '/logs?l='+filters.label;
      history.pushState(null, null, url);
    }
    app.func.stop(e);
  });
}

var TableRow = React.createClass({
  propTypes: {
    content: React.PropTypes.object.isRequired
  },
  statlog: function() {
    var tr = $(ReactDOM.findDOMNode(this)),
        id = tr.attr('data-container-id'),
        client = '?client='+tr.attr('data-client-id');
    app.func.link('/container/statlog/'+id+client);
  },
  render: function() {
    var log = this.props.content,
        key = log.id.substring(0, 20);
    return (
      <tr data-client-id={log.client} data-container-id={key}>
        <td className="data-name no-wrap">{log.time}</td>
        <td className="data-name"><a onClick={this.statlog}>{log.id.substring(0, 4)}</a></td>
        <td className="data-name">{log.type}</td>
        <td className="data-name data-force-break">{log.log.replace(/\[[0-9]+m/g, '')}</td>
      </tr>
    );
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
    var conditions = {count: monitoringCount};
    if (filters.client != -1) conditions.client = filters.client;

    app.func.ajax({url: '/api/logs', data: conditions, success: function (data) {
      if (data.error) {
        sender.setState({data: []});
        return;
      }
      stream = data;
      sender.setState({data: sender.filter()});
    }});
  },
  filter: function() {
    var data = [];

    // filter by labels
    var allowedIds = {};
    $.map(allcontainers, function (candidate) {
      $.map(candidate.containers, function (container) {
        if ((filters.label == 0) || ((filters.label == -1) && (! container.labels))) {
          allowedIds[container.id] = true;
        } else {
          if (! container.labels) return;
          var match = false;
          $.map(container.labels, function (value, key) {
            match |= (filters.label == app.func.hash(key+'->'+value));
          });
          if (match) allowedIds[container.id] = true;
        }
      });
    });

    // retrieve log data
    $.map(stream, function (host) {
      if ((filters.client != -1) && (host.client.id != filters.client)) return;
      $.map(host.logs, function (client) {
        if (! allowedIds[client.id]) return;

        $.map(client.stdout, function (record) {
          data.push({
            client: host.client.id,
            id: client.id,
            type: 'stdlog',
            key: record.substring(0, 30),
            time: record.substring(5, 19).replace('T', ' '),
            log: record.substring(31)
          });
        });
        $.map(client.stderr, function (record) {
          data.push({
            client: host.client.id,
            id: client.id,
            type: 'stderr',
            key: record.substring(0, 30),
            time: record.substring(5, 19).replace('T', ' '),
            log: record.substring(31)
          });
        });
      });
    });
    data.sort(function (a, b) {
      var diff = new Date(a.key.substring(0, 19)+'Z') - new Date(b.key.substring(0, 19)+'Z');
      if (diff != 0) return diff;
      return parseInt(a.key.substring(20), 10) - parseInt(b.key.substring(20), 10);
    })
    $('.logs').fadeIn();
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
    var rows = [];
    this.state.data.map(function(record) {
      if (! record.log) return;
      rows.push(<TableRow key={record.client + record.id + record.key} content={record} />)
    });
    return (
        <table className="table table-striped table-hover">
          <thead>
            <tr>
              <th>Time</th>
              <th>ID</th>
              <th>Type</th>
              <th>Log</th>
            </tr>
          </thead>
          <tbody>{rows}</tbody>
        </table>
    );
  }
});

ReactDOM.render(<Table reload={false} />, document.getElementById('data'));
