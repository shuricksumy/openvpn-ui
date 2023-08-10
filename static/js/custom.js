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
})

$(function() {
 $('a.btn-danger').confirm({
    content: "This cannot be reverted.",
    type: 'red',
    icon: 'fa fa-warning',
    title: 'font-awesome',
    columnClass: 'small',
    buttons: {
        Confirm: {
            btnClass: 'btn-red bg-red-gradient btn80',
            action: function(){location.href = this.$target.attr('href');}
        },
        Cancel: {
            btnClass: 'btn80',
            action: function(){}
        },
    }
   });

  $('a.btn-warning').confirm({
    content: "This will disable access for user",
    type: 'orange',
    icon: 'fa fa-warning',
    title: 'font-awesome',
    columnClass: 'small',
    buttons: {
        Confirm: {
            btnClass: 'btn-orange bg-yellow-gradient btn80',
            action: function(){location.href = this.$target.attr('href');}
        },
        Cancel: {
            btnClass: 'btn80',
            action: function(){}
        },
    }
   });
})

$(function() {
  $(document).ready(function () {
    $(".edit-button").on("click", function () {
        var clientName = $(this).data("client-name");
        showModalWithData(clientName);
    });
});

function showModalWithData(clientName) {
    // Make an AJAX request to get the data for the client
    $.post("/certificates/render_modal", {"client-name": clientName}, function (data) {
        // Update the modal content with the retrieved data
        $("#modal-edit-client").html(data);
        // Show the modal
        $("#editClientModal").modal("show");
      }).fail(function () {
        alert("Error loading data for the client.");
      }).always(function() {
        // Show the modal
        $("#editClientModal").modal("show");
    });
}
})

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


        // Copy the content to the clipboard when the Copy button is clicked
        $('.copy-btn').on('click', function () {
            var $popupContent = $(this).siblings('.param-name');
            var content = $popupContent.map(function () {
                return $(this).text() + $(this).next().text() + '\n';
            }).get().join('');

            // Create a temporary textarea to copy the content to the clipboard
            var $tempTextArea = $('<textarea>').val(content).css('position', 'absolute').css('left', '-9999px');
            $('body').append($tempTextArea);
            $tempTextArea.select();
            document.execCommand('copy');
            $tempTextArea.remove();
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