var isViewOnly = false;

$(document).ready(function () {
  $('#menu-clients').addClass('active');
  $('#client-detail pre').css({height: ($(window).height()-200)+'px'});
  isViewOnly = ($('#mode-view-only').val() == 'true');

  $('#add-client').on('show.bs.modal', function () {
    $('#endpoint').val('');
  });
  $('#add-client').on('shown.bs.modal', function () {
    $('#endpoint').focus();
  });
  $('#add-client .act-add').click(function () {
    var endpoint = app.func.trim($('#endpoint').val());
    if (endpoint.length == 0) {
      $('#endpoint').focus();
      return;
    }
    $('#add-client').modal('hide');
    _add(endpoint);
  });

  $('.container').bind('drop', function (e) {
    _uploadFile(e.originalEvent.dataTransfer.files);
    app.func.stop(e);
  }).bind('dragenter', function (e) {
    app.func.stop(e);
  }).bind('dragover', function (e) {
    app.func.stop(e);
  });
});

function _add(endpoint) {
  var arg = {endpoint: endpoint, cert: ''};
  app.func.ajax({type: 'POST', url: '/api/client/', data: arg, dataType: 'html', success: function () {
    ReactDOM.render(<Table />, document.getElementById('data'));
  }, error: function (xhr, status, err) {
    var message = (xhr && xhr.responseText) ? xhr.responseText : err;
    alert(message);
  }});
}

function _remove(id) {
  app.func.ajax({type: 'DELETE', url:'/api/client/'+id, dataType: 'html', success: function () {
    ReactDOM.render(<Table />, document.getElementById('data'));
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

function _uploadFile(files) {
  var fd = new FormData();
  for (var i = 0; i < files.length; i++) {
    fd.append("files[]", files[i]);
  }
  $.ajax({
    type: 'POST', url: '/clients/import', data: fd,
    processData: false, contentType: false,
    success: function () {
      ReactDOM.render(<Table />, document.getElementById('data'));
    }
  });
}

var TableRow = React.createClass({
  propTypes: {
    content: React.PropTypes.object.isRequired
  },
  handleDetail: function() {
    var tr = $(ReactDOM.findDOMNode(this)).closest('tr'),
        endpoint = tr.find('.endpoint').text(),
        cert = tr.find('.cert').text(),
        client = {endpoint: endpoint, cert: cert},
        popup = $('#client-detail');
    app.func.ajax({type: 'GET', url: '/api/client/', data: client, success: function (data) {
      popup.find('.detail-title').text(endpoint);
      popup.find('.details').text(JSON.stringify(data, true, ' '));
      popup.modal('show');
    }});
  },
  handleDelete: function() {
    var tr = $(ReactDOM.findDOMNode(this)).closest('tr'),
        id = tr.attr('data-client-id'),
        endpoint = tr.find('.endpoint').text();
    if (window.confirm('Are you sure to remove the client?\nEndpoint: '+endpoint)) {
      _remove(id);
    }
  },
  render: function() {
    var rowclass = this,
        client = this.props.content.client,
        info = this.props.content.info,
        version = this.props.content.version;
    return (
        <tr data-client-id={client.id}>
          <td className="data-index endpoint">
            <a onClick={rowclass.handleDetail} style={{outline: 'none', textDecoration: 'none'}}>{client.endpoint}</a>
          </td>
          <td className="data-index">{info.Containers}</td>
          <td className="data-index cert" style={{display: 'none'}}>{client.certPath}</td>
          <td className="data-index">{info.Name}</td>
          <td className="data-index">{_find(version, 'Version')}</td>
          <td className="data-index">{_find(version, 'ApiVersion')}</td>
          <td className="data-index">{client.isActive ? 'yes' : 'no'}</td>
          {isViewOnly ? '' : function() {
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
      return <TableRow key={index} content={record} />
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

ReactDOM.render(<Table />, document.getElementById('data'));
