var table;

$(document).ready(function () {
  $('#menu-clients').addClass('active');
  $('#client-detail pre').css({height: ($(window).height()-200)+'px'})

  $('#add-client').on('show.bs.modal', function (e) {
    $('#endpoint').val('');
  });
  $('#add-client').on('shown.bs.modal', function (e) {
    $('#endpoint').focus();
  });
  $('#add-client .act-add').click(function (e) {
    var endpoint = app.func.trim($('#endpoint').val());
    if (endpoint.length == 0) {
      $('#endpoint').focus();
      return;
    }
    $('#add-client').modal('hide');
    _add(endpoint);
  });
});

function _add(endpoint) {
  var arg = {endpoint: endpoint, cert: ''};
  app.func.ajax({type: 'POST', url:'/api/client/', data: arg, dataType: 'html', success: function (data) {
    table.setProps();
  }, error: function (xhr, status, err) {
    var message = (xhr && xhr.responseText) ? xhr.responseText : err;
    alert(message);
  }});
}

function _remove(id) {
  app.func.ajax({type: 'DELETE', url:'/api/client/'+id, dataType: 'html', success: function (data) {
    table.setProps();
  }});
}

function _find(arr, key, def) {
  var result = def ? def : '';
  if (! arr) {
    return result;
  }
  $.map(arr, function (value) {
    if (value.split('=')[0] == key) {
      result = value.split('=')[1];
    }
  });
  return result;
}

var TableRow = React.createClass({
  handleDetail: function() {
    var tr = $(this.getDOMNode()).closest('tr'),
        id = tr.attr('data-client-id'),
        endpoint = tr.find('.endpoint').text(),
        cert = tr.find('.cert').text(),
        client = {endpoint: endpoint, cert: cert},
        popup = $('#client-detail');
    app.func.ajax({type: 'POST', url: '/api/client/', data: client, success: function (data) {
      popup.find('.detail-title').text(endpoint);
      popup.find('.details').text(JSON.stringify(data, true, ' '));
      popup.modal('show');
    }});
    return false;
  },
  handleDelete: function() {
    var tr = $(this.getDOMNode()).closest('tr'),
        id = tr.attr('data-client-id'),
        endpoint = tr.find('.endpoint').text();
    if (window.confirm('Are you sure to remove the client?\nEndpoint: '+endpoint)) {
      _remove(id);
    }
    return false;
  },
  render: function() {
    var rowclass = this,
        client = this.props.content.client,
        info = this.props.content.info,
        version = this.props.content.version;
    return (
        <tr key={this.props.index} data-client-id={client.id}>
          <td className="data-index endpoint">
            <a href="#" onClick={rowclass.handleDetail} style={{outline: 'none', textDecoration: 'none'}}>{client.endpoint}</a>
          </td>
          <td className="data-index">{_find(info, 'Containers')}</td>
          <td className="data-index cert" style={{display: 'none'}}>{client.certPath}</td>
          <td className="data-index">{_find(info, 'Name')}</td>
          <td className="data-index">{_find(version, 'Version')}</td>
          <td className="data-index">{_find(version, 'ApiVersion')}</td>
          <td className="data-index">{client.isActive ? 'yes' : 'no'}</td>
          {client.isDefault ? '' : function() {
            return (
              <td style={{width: 60+'px'}}>
                <a className="btn btn-danger" onClick={rowclass.handleDelete}>&times;</a>
              </td>
            )
          }()}
        </tr>
    );
  }
});

var Table = React.createClass({
  getInitialState: function() {
    return {data: []};
  },
  load: function(sender) {
    app.func.ajax({type: 'GET', url: '/api/clients', success: function (data) {
      $('#count').text(data.length + ' client' + ((data.length > 1) ? 's' : ''));
      data.sort(function (a, b) {
        return parseInt(_find(a.info, 'Name'), 10) - parseInt(_find(b.info, 'Name'), 10);
      });
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
      return <TableRow index={index} content={record} />
    });
    return (
        <table className="table table-striped table-hover">
          <thead>
            <tr>
              <th>Endpoint</th>
              <th>Containers</th>
              <th>Host Name</th>
              <th>Version</th>
              <th>API Version</th>
              <th>Active</th>
              <th></th>
            </tr>
          </thead>
          <tbody>{rows}</tbody>
        </table>
    );
  }
});

table = React.render(<Table />, document.getElementById('data'));
