$.MyAPP = {};

$.MyAPP.Disconnect = function (cname){
  console.log(cname)
  $.ajax({
    type: "DELETE",
    dataType: "json",
    url: "api/v1/session",
    data: JSON.stringify({ "cname": cname }),
    success: function(data) {
      location.reload();
      console.log(data);
    },
    error: function(a,b,c) {
      console.log(a,b,c)
      location.reload();
    }
  });
}

$(function() {
  new Clipboard('.button-copy');

  //$( ".btn-disconnect" ).click(function() {
  //  alert( "Handler for .click() called." );
  //});
  //window.location.reload();
});

$(function() {
 $('a.btn-danger').confirm({
    content: "Are you sure? This cannot be reverted.",
    type: 'red',
    icon: 'fa fa-hand-o-right',
    title: 'font-awesome',
    theme: 'bootstrap',
    columnClass: 'small',
    animateFromElement: false,
    animation: 'RotateXR',
    closeAnimation: 'rotateXR',
    buttons: {
        Confirm: {
            btnClass: 'btn-red bg-red-gradient btn80 myspiner',
            action: function(){location.href = this.$target.attr('href');}
        },
        Cancel: {
            btnClass: 'btn80',
            action: function(){}
        },
    }
   });

  $('a.btn-warning').confirm({
    content: "Confirm, if you are 100% sure.",
    type: 'orange',
    icon: 'fa fa-hand-o-right',
    title: 'font-awesome',
    columnClass: 'small',
    theme: 'bootstrap',
    animateFromElement: false,
    animation: 'RotateXR',
    closeAnimation: 'rotateXR',
    buttons: {
        Confirm: {
            btnClass: 'btn-orange bg-yellow-gradient btn80 setloader myspiner',
            action: function(){location.href = this.$target.attr('href');}
        },
        Cancel: {
            btnClass: 'btn80',
            action: function(){}
        },
    }
   });

   $('button.form-to-confirm').confirm({
    content: "Confirm, if you are 100% sure.",
    type: 'blue',
    icon: 'fa fa-hand-o-right',
    title: 'font-awesome',
    columnClass: 'small',
    theme: 'bootstrap',
    animateFromElement: false,
    animation: 'RotateXR',
    closeAnimation: 'rotateXR',
    buttons: {
        Confirm: {
            btnClass: 'btn-primary bg-light-blue-gradient btn80 setloader myspiner',
            action: function(){  $('#form-to-confirm').submit(); }
        },
        Cancel: {
            btnClass: 'btn80',
            action: function(){}
        },
    }
   });
})


// FUNCTION EDIT CLIENT DETAILS POPUP RAW TEXT EDITOR
$(function() {
    $(document).ready(function () {
        $("button#openModalEditClientRaw").on("click", function () {
            var clientName = $(this).data("client-name");
            showModalWithData(clientName);
        });
    });

    function showModalWithData(clientName) {
        // Make an AJAX request to get the data for the client
        $.post("/clients/render_modal_raw", {"client-name": clientName}, function (data) {
            // Update the modal content with the retrieved data
            $("#modal-edit-client-raw").html(data);
            // Show the modal
            $("#editClientModalRaw").modal("show");
        }).fail(function () {
            alert("Error loading data for the client.");
        }).always(function() {
            // Show the modal
            $("#editClientModalRaw").modal("show");
        });
    }
});

// FUNCTION EDIT CLIENT DETAILS POPUP FORM WITH ROUTES
$(function() {
    $(document).ready(function () {
        $("button#openModalEditClientDetails").on("click", function () {
            var clientName = $(this).data("client-name");
            showModalWithDataClientDetials(clientName);
        });
    }); $()

    function showModalWithDataClientDetials(clientName) {
        // Make an AJAX request to get the data for the client
        $.post("/clients/render_modal", {"client-name": clientName}, function (data) {
            // Update the modal content with the retrieved data
            $("#modal-edit-client-details").html(data);
            // Show the modal
            $("#editClientDetailsModal").modal("show");
        }).fail(function () {
            alert("Error loading data for the client.");
        }).always(function() {
            // Show the modal
            $("#editClienDetailstModal").modal("show");
        });
    }
});

$(function() {
    $(document).ready(function () {
        $("button#getRouteEditModal").on("click", function () {
            var routeID = $(this).data("route-id");
            showEditRouteModal(routeID);
        });
    }); $()

    function showEditRouteModal(routeID) {
        // Make an AJAX request to get the data for the client
        $.get("/routes/get/" + routeID, function (data) {
            // Update the modal content with the retrieved data
            $("#modal-edit-route-details").html(data);
            // Show the modal
            $("#showEditRouteModal").modal("show");
        }).fail(function () {
            alert("Error loading data for the client.");
        }).always(function() {
            // Show the modal
            $("#showEditRouteModal").modal("show");
        });
    }
});

$(document).ready(function() {
    // Use $(document).on() to bind the click event to dynamically added elements
    $(document).on('click', '[data-popup-target]', function() {
        var targetId = $(this).data('popup-target');
        var targetPopup = $('#' + targetId);

        // Show the popup
        targetPopup.show();

        // Hide the popup after 5 seconds
        setTimeout(function() {
            targetPopup.hide();
        }, 5000);
    });
});

function createEditor(name, size, theme, mode, readonly) {
    // find the textarea
    var textarea = document.querySelector("form textarea[name=" + name + "]");

    // create ace editor 
    var editor = ace.edit()
    editor.container.style.height = size

    editor.setTheme("ace/theme/" + theme); //"clouds_midnight"
    //editor.setTheme("ace/theme/twilight");
    //editor.setTheme("ace/theme/iplastic");

    editor.session.setMode("ace/mode/" + mode);

    editor.setReadOnly(readonly);
    editor.setShowPrintMargin(false);
    editor.session.setUseWrapMode(true);
    editor.session.setValue(textarea.value)
    // replace textarea with ace
    textarea.parentNode.insertBefore(editor.container, textarea)
    textarea.style.display = "none"
    // find the parent form and add submit event listener
    var form = textarea
    while (form && form.localName != "form") form = form.parentNode
    form.addEventListener("submit", function() {
        // update value of textarea to match value in ace
        textarea.value = editor.getValue()
    }, true)
}

$(".reveal").on('click',function() {
    var $pwd = $(".pwd");
    if ($pwd.attr('type') === 'password') {
        $pwd.attr('type', 'text');
    } else {
        $pwd.attr('type', 'password');
    }
});

jQuery(function(){


    $(document).on("click", ".myspiner", function() {
        $("#overlay").fadeIn(300);
    });

  
  $('.myspiner').click(function(){
    $.ajax({
      type: 'GET',
      success: function(data){
        console.log(data);
      }
    }).done(function() {
      setTimeout(function(){
        $("#overlay").fadeOut(300);
      },500);
    });
  });	
});



$(document).ready(function() {
    $('.select2').select2();
});

  function updateStatus() {
      $.ajax({
          url: "/openvpn/status",
          type: "GET",
          success: function(response) {
              if (response.includes("running")) {
                  $("#openvpn-status").removeClass("stopped").addClass("running").attr("title", "OpenVPN is running");
              } else {
                  $("#openvpn-status").removeClass("running").addClass("stopped").attr("title", "OpenVPN is stopped");
              }
          },
          error: function() {
              $("#openvpn-status").removeClass("running").addClass("stopped").attr("title", "OpenVPN is stopped");
          }
      });
  }
  
  function startOpenVPN() {
      $.ajax({
          url: "/openvpn/start",
          type: "GET",
          success: function(response) {
              updateStatus();
              $("#start-stop-messages").css("color", "green");
              $("#start-stop-messages").text(response);
          },
          error: function(err) {
              $("#start-stop-messages").css("color", "red");
              $("#start-stop-messages").text("Error starting OpenVPN: " + err.responseText);
          }
      });
  }
  
  function stopOpenVPN() {
      $.ajax({
          url: "/openvpn/stop",
          type: "GET",
          success: function(response) {
              updateStatus();
              $("#start-stop-messages").css("color", "green");
              $("#start-stop-messages").text(response);
          },
          error: function(err) {
              $("#start-stop-messages").css("color", "red");
              $("#start-stop-messages").text("Error stopping OpenVPN: " + err.responseText);
          }
      });
  }
  
  // Update status every minute
  setInterval(updateStatus, 60000);
  
  // Initial update
  updateStatus();
