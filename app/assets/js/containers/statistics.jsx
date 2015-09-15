var table = false,
    statistics = {
      previous: [],
      current: []
    },
    storedCPUs = [],
    storedMems = [],
    refreshWindow = app.storage.get('refresh-statistics', 2),
    clients = 0;

$(document).ready(function () {
  $('#menu-statistics').addClass('active');
  clients = parseInt($('#number-of-clients').val(), 10);

  setRefreshWindow(refreshWindow);
  $('#refresh-window a').click(function(e) {
    setRefreshWindow(parseInt($(this).attr('href').substring(1), 10));
    app.func.stop(e);
  });
  refreshStats();
});

function setRefreshWindow(value) {
  var a = $('#refresh-window a[href="#'+value+'"]'),
      group = a.closest('.btn-group').removeClass('open');
  refreshWindow = value;
  app.storage.set('refresh-statistics', value);
  group.find('.caption').text('refresh / '+a.text()).blur();
}

function refreshStats() {
  if (refreshWindow && (refreshWindow > 0)) {
    table && table.setProps();
  }
  setTimeout(refreshStats, Math.max(1, refreshWindow) * 1000);
}

var lineChartCPU = nv.models.lineChart()
              .x(function(d) { return d[0] })
              .y(function(d) { return d[1]/100 })
              .color(d3.scale.category10().range())
              .showLegend(false)
              .useInteractiveGuideline(true);
lineChartCPU.xAxis.tickFormat(function(d) {return d3.time.format('%X')(new Date(d))});
lineChartCPU.yAxis.tickFormat(d3.format(',.1%'));
nv.utils.windowResize(lineChartCPU.update);

var lineChartMem = nv.models.lineChart()
              .x(function(d) { return d[0] })
              .y(function(d) { return d[1]/100 })
              .color(d3.scale.category10().range())
              .showLegend(false)
              .useInteractiveGuideline(true);
lineChartMem.xAxis.tickFormat(function(d) {return d3.time.format('%X')(new Date(d))});
lineChartMem.yAxis.tickFormat(d3.format(',.1%'));
nv.utils.windowResize(lineChartMem.update);

if (clients <= 1) {
  var pieChartCPU = nv.models.pieChart()
              .x(function(d) { return d.label })
              .y(function(d) { return d.value })
              .color(d3.scale.category10().range())
              .showLegend(false)
              .showLabels(true)
              .labelThreshold(.05)
              .labelType("percent")
              .donut(true)
              .donutRatio(0.35);

  var pieChartMem = nv.models.pieChart()
              .x(function(d) { return d.label })
              .y(function(d) { return d.value })
              .color(d3.scale.category10().range())
              .showLegend(false)
              .showLabels(true)
              .labelThreshold(.05)
              .labelType("percent")
              .donut(true)
              .donutRatio(0.35);
}

function _setStoredData(arr, key, values) {
  var found = false;
  $.map(arr, function (record, index) {
    if (record.key == key) found = record;
  });
  if (found) {
    found.values.push(values);
    if (found.values.length > 20) found.values.shift()
    return;
  }
  arr.push({key: key, values: [values]});
}

var TableRow = React.createClass({
  render: function() {
    var name = this.props.content.name,
        stat = this.props.content.current && this.props.content.current[0],
        prev = this.props.content.previous && this.props.content.previous[0],
        cpu_delta = 0, system_delta = 0, cpu_percent = 0;
    if (prev && stat && stat.cpu_stats) {
      cpu_delta = stat.cpu_stats.cpu_usage.total_usage - prev.cpu_stats.cpu_usage.total_usage;
      system_delta = stat.cpu_stats.system_cpu_usage - prev.cpu_stats.system_cpu_usage;
    }
    if ((system_delta > 0) && (cpu_delta > 0)) {
      cpu_percent = 100.0 * cpu_delta / system_delta * stat.cpu_stats.cpu_usage.percpu_usage.length;
    }
    var time = '',
        mem = {usage: '-', max: '-', limit: '-', percent: 0},
        network = {in: '', out: '', inPacket: '', outPacket: ''};
    if (stat && stat.read) {
      time = stat.read.substring(5, 19).replace(/-/, '/').replace('T', ' ');
      mem = {
        usage: app.func.byteFormat(stat.memory_stats.usage),
        max: app.func.byteFormat(stat.memory_stats.max_usage),
        percent: stat.memory_stats.usage * 100 / stat.memory_stats.limit
      };
      network = {
        in: app.func.byteFormat(stat.network.rx_bytes),
        out: app.func.byteFormat(stat.network.tx_bytes)
      };
    }
    return (
        <tr key={this.props.index}>
          <td className="data-name">{name.substring(1).replace(',/', ', ')}</td>
          <td className="data-name">{time}</td>
          <td className="data-name">{(cpu_percent+'').substring(0, 4)}%</td>
          <td className="data-name">{mem.usage} / {mem.max}</td>
          <td className="data-name">{(mem.percent+'').substring(0, 4)}%</td>
          <td className="data-name">{network.in} / {network.out}</td>
        </tr>
    );
  }
});

var Table = React.createClass({
  getInitialState: function() {
    return {data: {previous: [], current: []}};
  },
  load: function(sender) {
    var self = this;
    app.func.ajax({type: 'GET', url: '/api/statistics', success: function (data) {
      data.sort(function (a, b) {
        if (a.client.endpoint < b.client.endpoint)
          return -1;
        if (a.client.endpoint > b.client.endpoint)
          return 1;
        return 0;
      });
      if (data.error) {
        statistics.previous = true;
        statistics.current = data;
      } else {
        statistics.previous = statistics.current;
        statistics.current = data;
      }
      // update stats-table
      sender.setState({data: statistics});

      var stats = [];
      $.map(data, function (host) {
        $.map(host.stats, function (record, key) {
          stats.push({
            client: host.client, key: key,
            stat: record && record[0]
          })
        });
      });
      stats.sort(function (a, b) {
        var an = a.key + a.client.endpoint,
            bn = b.key + b.client.endpoint;
        if (an < bn) return -1;
        if (an > bn) return  1;
        return 0;
      });

      // change data format for charts
      var pie = {CPU: [], Mem: []}, multiple = (data.length > 1);
      $.map(stats, function (record, index) {
        var client = record.client,
            name = record.key.substring(1).replace(',/', ', ') + _endpoint(multiple, client.endpoint),
            stat = record.stat,
            time = new Date(stat.read.substring(0, 19)+'Z').getTime(),
            prev = _findPrivious(client.endpoint, record.key),
            cpu_delta = 0, system_delta = 0, cpu_percent = 0, mem_percent = 0;
        if (prev && (prev.length > 0)) {
          prev = prev[0];
        }
        if (prev && stat && stat.cpu_stats) {
          cpu_delta = stat.cpu_stats.cpu_usage.total_usage - prev.cpu_stats.cpu_usage.total_usage;
          system_delta = stat.cpu_stats.system_cpu_usage - prev.cpu_stats.system_cpu_usage;
        }
        if ((system_delta > 0) && (cpu_delta > 0)) {
          cpu_percent = 100.0 * cpu_delta / system_delta * stat.cpu_stats.cpu_usage.percpu_usage.length;
        }
        if (stat) mem_percent = stat.memory_stats.usage * 100 / stat.memory_stats.limit;

        _setStoredData(storedCPUs, name, [time, cpu_percent]);
        pie.CPU.push({label: name, value: cpu_percent});
        _setStoredData(storedMems, name, [time, mem_percent]);
        pie.Mem.push({label: name, value: mem_percent});
      });

      // draw line-charts
      nv.addGraph(function() {
        d3.select('#chart-cpu svg.line-charts').datum(storedCPUs).call(lineChartCPU);
        var max = d3.max(storedCPUs, function(d) {
          return d3.max(d.values, function(value) {
            return (value[1] + 0.1) / 100;
          });
        });
        return lineChartCPU.yDomain([0, max]);
      });
      if (clients <= 1) {
        var remain = 100;
        $.map(pie.CPU, function (record) {
          remain -= record.value;
        });
        pie.CPU.push({label: '-', value: remain});
        d3.select("#chart-cpu svg.pie-charts").datum(pie.CPU).transition().duration(350).call(pieChartCPU);
      }
      nv.addGraph(function() {
        d3.select('#chart-mem svg.line-charts').datum(storedMems).call(lineChartMem);
        var max = d3.max(storedMems, function(d) {
          return (d.values[d.values.length-1][1] + 0.5) / 100;
        });
        return lineChartMem.yDomain([0, max]);
      });
      if (clients <= 1) {
        var remain = 100;
        $.map(pie.Mem, function (record) {
          remain -= record.value;
        });
        pie.Mem.push({label: '-', value: remain});
        d3.select("#chart-mem svg.pie-charts").datum(pie.Mem).transition().duration(350).call(pieChartMem);
      }
    }});
  },
  componentDidMount: function() {
    this.load(this);
  },
  componentWillReceiveProps: function() {
    this.load(this);
  },
  render: function() {
    var data = this.state.data, rows = [], multiple = (data.current.length > 1);
    $.map(data.current, function (host, index) {
      var client = host.client;

      $.map(host.stats, function (current, key) {
        var name = key + _endpoint(multiple, client.endpoint);
        rows.push(<TableRow key={key+'@'+client.id} index={key+'@'+index} content={{
          name: name,
          previous: _findPrivious(client.endpoint, key),
          current: current
        }} />)
      });
    });
    return (
        <table className="table table-striped table-hover">
          <thead>
            <tr>
              <th>ID</th>
              <th>Time</th>
              <th>CPU %</th>
              <th>MEM USAGE / MAX</th>
              <th>MEM %</th>
              <th>NET I/O</th>
            </tr>
          </thead>
          <tbody>{rows}</tbody>
        </table>
    );
  }
});

function _endpoint(multiple, endpoint) {
  return (multiple ? ' @'+endpoint.replace(/^.*:\/\//, '').replace(/:.*$/, '') : '');
}

function _findPrivious(endpoint, key) {
  var result = false;
  if (statistics.previous) {
    $.map(statistics.previous, function (host) {
      if (host.client.endpoint == endpoint)
        result = host.stats[key];
    });
  }
  return result;
}

table = React.render(<Table />, document.getElementById('data'));
