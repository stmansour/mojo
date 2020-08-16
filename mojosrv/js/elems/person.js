"use strict";
/*global
    w2ui, $, app, console, w2utils, document
*/

var addGroup = {set: false, GroupName: '', GID: 0};

function resetAddGroup() {
    addGroup.set = false;
    addGroup.GroupName = '';
    addGroup.GID = 0;
}

function setAddGroup(name,id) {
    addGroup.set = true;
    addGroup.GroupName = name;
    addGroup.GID = id;
}

//------------------------------------------------------------------------
//          Person Form
//------------------------------------------------------------------------
$().w2form({
    name: 'personForm',
    style: 'border: 0px; background-color: transparent;',
    header: 'Person Detail',
    url: '/v1/person',
    formURL: '/html/formperson.html',
    fields: [
        { field: 'recid',          type: 'int',   required: false },
        { field: 'PID',            type: 'text',  required: false },
        { field: 'FirstName',      type: 'text',  required: false },
        { field: 'MiddleName',     type: 'text',  required: false },
        { field: 'LastName',       type: 'text',  required: false },
        { field: 'PreferredName',  type: 'text',  required: false },
        { field: 'JobTitle',       type: 'text',  required: false },
        { field: 'OfficePhone',    type: 'text',  required: false },
        { field: 'OfficeFax',      type: 'text',  required: false },
        { field: 'Email1',         type: 'email', required: true },
        { field: 'Email2',         type: 'email', required: false },
        { field: 'MailAddress',    type: 'text',  required: false },
        { field: 'MailAddress2',   type: 'text',  required: false },
        { field: 'MailCity',       type: 'text',  required: false },
        { field: 'MailState',      type: 'text',  required: false },
        { field: 'MailPostalCode', type: 'text',  required: false },
        { field: 'MailCountry',    type: 'text',  required: false },
        { field: 'RoomNumber',     type: 'text',  required: false },
        { field: 'MailStop',       type: 'text',  required: false },
        { field: 'Status',         type: 'int',   required: false },
        { field: 'OptOutDate',     type: 'date',  required: false },
        { field: 'LastModTime',    type: 'text',  required: false },
        { field: 'LastModBy',      type: 'text',  required: false },
        { field: 'GroupName',      type: 'enum',  required: false,
            options: {
                url:            '/v1/grouptd/',
                max:            3,
                items:          [],
                openOnFocus:    false,
                maxDropWidth:   500,
                maxDropHeight:  500,
                renderItem:     groupPickerRender,
                renderDrop:     groupPickerDropRender,
                compare:        groupPickerCompare,
                onNew: function (event) {
                    $.extend(event.item, { GID: 0, GroupName: '' } );
                },
                onRemove: function(event) {
                    // event.onComplete = function() {
                    //     var f = w2ui.bizDetailForm;
                    //     // reset BUD field related data when removed
                    //     f.record.ClassCode = 0;
                    //     f.record.CoCode = 0;
                    //     f.record.BUD = "";
                    //
                    //     // NOTE: have to trigger manually, b'coz we manually change the record,
                    //     // otherwise it triggers the change event but it won't get change (Object: {})
                    //     var event = f.trigger({ phase: 'before', target: f.name, type: 'change', event: event }); // event before
                    //     if (event.cancelled === true) return false;
                    //     f.trigger($.extend(event, { phase: 'after' })); // event after
                    // };
                },
            },
        },

    ],
    toolbar: {
        items: [
            { id: 'btnNotes', type: 'button', icon: 'fa fa-sticky-note-o' },
            { id: 'bt3', type: 'spacer' },
            { id: 'btnClose', type: 'button', icon: 'fa fa-times' },
        ],
        onClick: function (event) {
            if (event.target == 'btnClose') {
                        w2ui.toplayout.hide('right',true);
                        w2ui.peopleGrid.render();
            }
        },
    },
});

//------------------------------------------------------------------------
//    personGroupsGrid
//
//                top    = resUpdateForm
//    >>>>>>>>>   main   = personGroupsGrid    <<<<<<<<<<<
//                bottom = personFormBtns
//------------------------------------------------------------------------
$().w2grid({
    name: 'personGroupsGrid',
    url: '/v1/pgroup',
    header: 'Group Membership',
    multiSelect: false,
    show: {
        toolbar         : true,
        footer          : false,
        header          : false,
        toolbarAdd      : false,   // indicates if toolbar add new button is visible
        toolbarDelete   : true,    // indicates if toolbar delete button is visible
        toolbarSave     : false,   // indicates if toolbar save button is visible
        selectColumn    : false,
        expandColumn    : false,
        toolbarEdit     : false,
        toolbarSearch   : false,
        toolbarInput    : false,
        searchAll       : false,
        toolbarReload   : false,
        toolbarColumns  : false,
    },
    columns: [
        {field: 'recid',     caption: 'recid', size: '40px',  hidden: true,  sortable: true },
        {field: 'GID',       caption: 'GID',   size: '50px',  hidden: false, sortable: true },
        {field: 'GroupName', caption: 'Group', size: '350px', hidden: false, sortable: true },
    ],
    onLoad: function(event) {
        event.onComplete = function() {
            w2ui.personGroupsGrid.url = '';  // no updates at this point or we will wipe out any changes the person makes
        };
    },
    onDelete: function(event) {
        event.onComplete = function() {
            // console.log('Delete');
        };
    },
    onClick: function(event) {
        event.onComplete = function () {
            var rec = w2ui.personGroupsGrid.get(event.recid);
        };
    },
    onRefresh: function(event) {
        event.onComplete = function() {
        };
    }
});


//------------------------------------------------------------------------
//    personFormBtns
//
//                top    = resUpdateForm
//                main   = personGroupsGrid
//    >>>>>>>>>   bottom = personFormBtns      <<<<<<<<<
//------------------------------------------------------------------------
$().w2form({
    name: 'personFormBtns',
    style: 'border: 0px; background-color: transparent;',
    formURL: '/html/formpersonbtns.html',
    url: '',
    fields: [],
    actions: {
        save: function () {

            $.when(
                savePersonForm(),
                savePersonGroups(),
            )
            .done(function(){
                personSaveDoneCB();
            })
            .fail(function(){
                var s = 'Person Save: error reported';
                w2ui.personGrid.error(s);
                personSaveDoneCB();
            });



        },
        saveadd: function() {

        },
        delete: function() {
            var request={cmd:"delete",selected: [w2ui.personForm.record.PID]};
            $.post('/v1/person/'+w2ui.personForm.record.PID, JSON.stringify(request))
            .done(function(data) {
                if (typeof data == 'string') {  // it's weird, a successful data add gets parsed as an object, an error message does not
                    var msg = JSON.parse(data);
                    w2ui.personForm.error(msg.message);
                    return;
                }
                if (data.status != 'success') {
                    w2ui.personForm.error(data.message);
                }
            });
            w2ui.toplayout.hide('right',true);
            w2ui.peopleGrid.reload();
        },
    },
});


//------------------------------------------------------------------------
//    savePersonForm saves the data in the top form of the Person UI.
//------------------------------------------------------------------------
function savePersonForm() {
    var rec = w2ui.personForm.record;
    var params = {
        cmd: "save",
        selected: [],
        limit: 0,
        offset: 0,
        record: rec,
    };
    var dat = JSON.stringify(params);
    return $.post(w2ui.personForm.url, dat, null, "json")
    .done(function (data) {
        if (data.status === 'error') {
            console.log('ERROR: ' + data.message);
            return;
        }
    })
    .fail( function() {
            w2ui.rentablesGrid.error('post failed to ' + w2ui.rentableForm.url);
        }
    );

}

// savePersonGroups looks at all the rows in the person's group grid and sends
// back to the server the list of IDs that make up the person's group
// membership. The server will update accordingly.
//
// Format expected:  /v1/groupmembership/PID
//         payload:  {"cmd":"save","Groups":[4,3,2,5]}
//
function savePersonGroups() {
    var g = w2ui.personGroupsGrid;
    var r = g.records;
    var PID = w2ui.personForm.record.PID;
    var url = '/v1/groupmembership/' + PID;
    var rec = {
        cmd: "save",
        Groups: [],
    };
    for (var i = 0; i < r.length; i++) {
        rec.Groups.push( r[i].GID);
    }
    var dat = JSON.stringify(rec);

    return $.post(url, dat, null, "json")
    .done(function(data) {
        if (data.status === "success") {
            //------------------------------------------------------------------
            // Now that the save is complete, we can add the URL back to the
            // the grid so it can call the server to get updated rows. The
            // onLoad handler will reset the url to '' after the load completes
            // so that changes are done locally to gthe grid until the
            // rentableForm save button is clicked.
            //------------------------------------------------------------------
            w2ui.personGroupsGrid.url = '/v1/pgroup';
        } else {
            w2ui.personGroupsGrid.error('Save Groups: '+data.message);
        }
    })
    .fail(function(data){
        w2ui.personGroupsGrid.error("Save RentableUseStatus failed. " + data);
    });
}

// personSaveDoneCB handles removing the personDetail form from the page and
// reloading the people grid.
//------------------------------------------------------------------------------
function personSaveDoneCB() {
    w2ui.toplayout.hide('right',true);
    w2ui.peopleGrid.reload();
}

// addGroupHandler is called when the user clics the Add button to add a group.
//------------------------------------------------------------------------------
function addGroupHandler() {
    if ( typeof w2ui.personForm.record === "undefined" ) {return;}
    if ( typeof w2ui.personForm.record.GroupName === "undefined" ) {return;}
    if ( w2ui.personForm.record.GroupName === null) {return;}

    var g = w2ui.personGroupsGrid;
    var r = g.records;

    if (addGroup.set) {
        // if we already have this group, then just return
        for (var i = 0; i < r.length; i++) {
            if ( r[i].GID === addGroup.GID) {
                w2ui.personForm.record.GroupName = {};
                w2ui.personForm.refresh();
                return;
            }
        }
        // if we hit this point, then we need to add the group...
        var rec = {
            recid: addGroup.GID,
            GID: addGroup.GID,
            GroupName: addGroup.GroupName,
        };
        g.add(rec);
    }
    w2ui.personForm.record.GroupName = {};
    w2ui.personForm.refresh();
}

//-----------------------------------------------------------------------------
// groupPickerCompare - Compare item to the search string. Verify that the
//          supplied search string can be found in item.  If it's already
//          listed we don't want to list it again.
// @params
//   item = an object assumed to have a Name and GID field
// @return - true if the search string is found, false otherwise
//-----------------------------------------------------------------------------
function groupPickerCompare (item, search) {
    var s = item.GroupName;
    s = s.toLowerCase();
    var srch = search.toLowerCase();
    var match = (s.indexOf(srch) >= 0);
    return match;
}

//-----------------------------------------------------------------------------
// groupPickerDropRender - renders a name during typedown.
// @params
//   item = an object assumed to have a FirstName and LastName
// @return - the name to render
//-----------------------------------------------------------------------------
function groupPickerDropRender (item) {
    return item.GroupName;
}

//-----------------------------------------------------------------------------
// groupPickerRender - renders a name during typedown.
// @params
//   item = an object assumed to have a FirstName and LastName
// @return - true if the names match, false otherwise
//-----------------------------------------------------------------------------
function groupPickerRender (item) {
    setAddGroup(item.GroupName,item.GID);
    return item.GroupName;
}
