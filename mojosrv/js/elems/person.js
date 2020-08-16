"use strict";
"esversion: 8";
/*global
    w2ui, $, app, console, w2utils, document,
    getStatusString, showForm, initializePersonRecord,

*/

/*

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

function createPersonUI() {
    //------------------------------------------------------------------------
    //          peopleGrid
    //------------------------------------------------------------------------
    $().w2grid({
        name: 'peopleGrid',
        url: '/v1/people',
        show: {
            toolbar         : true,
            footer          : true,
            toolbarAdd      : true,   // indicates if toolbar add new button is visible
            toolbarDelete   : false,   // indicates if toolbar delete button is visible
            toolbarSave     : false,   // indicates if toolbar save button is visible
            selectColumn    : false,
            expandColumn    : false,
            toolbarEdit     : false,
            toolbarSearch   : false,
            toolbarInput    : true,
            searchAll       : false,
            toolbarReload   : true,
            toolbarColumns  : true,
        },
         columns: [
            {field: 'PID',          caption: 'PID',           size:  '60px', sortable: true, hidden: true},
            {field: 'Status',       caption: 'Status',        size: '100px', sortable: true, hidden: false,
                    render: function (record, index, col_index) {
                                return getStatusString(this.getCellValue(index, col_index));
                            }
            },
            {field: 'OptOutDate',   caption: 'OptOutDate',    size:  '80px', sortable: true, hidden: false},
            {field: 'Email1',       caption: 'Email1',        size: '200px', sortable: true, hidden: false},
            {field: 'Email2',       caption: 'Email2',        size:  '40px', sortable: true, hidden: false},
            {field: 'FirstName',    caption: 'FirstName',     size: '100px', sortable: true, hidden: false},
            {field: 'MiddleName',   caption: 'MiddleName',    size: '50px',  sortable: true, hidden: false},
            {field: 'LastName',     caption: 'LastName',      size: '100px', sortable: true, hidden: false},
            {field: 'PreferredName',caption: 'PreferredName', size: '100px', sortable: true, hidden: true},
            {field: 'JobTitle',     caption: 'JobTitle',      size: '150px', sortable: true, hidden: false},
            {field: 'OfficePhone',  caption: 'OfficePhone',   size: '100px', sortable: true, hidden: false},
            {field: 'OfficeFax',    caption: 'OfficeFax',     size: '100px', sortable: true, hidden: true},
            {field: 'MailAddress',  caption: 'MailAddress',   size: '100px', sortable: true, hidden: false},
            {field: 'MailAddress2', caption: 'MailAddress2',  size: '100px', sortable: true, hidden: true},
            {field: 'MailCity',     caption: 'MailCity',      size: '100px', sortable: true, hidden: false},
            {field: 'MailState',    caption: 'MailState',     size:  '50px', sortable: true, hidden: false},
            {field: 'MailPostalCode',caption:'MailPostalCode',size:  '50px', sortable: true, hidden: false},
            {field: 'MailCountry',  caption: 'MailCountry',   size: '100px', sortable: true, hidden: true},
            {field: 'RoomNumber',   caption: 'RoomNumber',    size: '100px', sortable: true, hidden: false},
            {field: 'MailStop',     caption: 'MailStop',      size: '100px', sortable: true, hidden: false},
            {field: 'LastModTime',  caption: 'LastModTime',   size: '100px', sortable: true, hidden: true},
        ],
        onClick: function(event) {
            event.onComplete = function (event) {
                var rec = w2ui.peopleGrid.get(event.recid);
                //w2ui.personForm.record = initializePersonRecord();
                w2ui.personForm.recid = rec.PID;
                w2ui.personForm.url = '/v1/person/' + rec.PID;
                w2ui.personForm.reload();
                w2ui.personGroupsGrid.url = '/v1/pgroup/' + rec.PID;
                w2ui.personGroupsGrid.reload();
                showForm(w2ui.personFormLayout, 500);
            };
        },
        onRequest: function(event) {
            w2ui.peopleGrid.postData = {groupName: app.groupFilter};
        },
        onAdd: function (/*event*/) {
            // Always create epoch assessment
            var f = w2ui.personForm;
            f.record = initializePersonRecord();
            f.recid = 0;
            f.refresh();
            w2ui.personGroupsGrid.url = '';  // no group associations at this time
            w2ui.personGroupsGrid.clear(); // remove anything that might be in there
            showForm(w2ui.personFormLayout,500);
        },
        onRefresh: function(/*event*/) {
            document.getElementById('mojoGroupFilter').value = app.groupFilter;
        },
        onLoad: function(/*event*/) {
            document.getElementById('mojoGroupFilter').value = app.groupFilter;
        },
        onSearch: function(event) {
            console.log('onSearch event fired. event = ' + event);
        }
    });

    w2ui.peopleGrid.toolbar.add(
        [
            { type: 'break', id: 'break1' },
            { type: 'html', id: 'groupName',
                    html: 'Group: <input type="text" id="mojoGroupFilter" value=""' +
                    'onkeypress="if (event.keyCode == 13) groupFilterSubmit();"' +
                    'value="'+app.groupFilter+'"'+
                    '>' },
        ],
    );

    w2ui.peopleGrid.toolbar.on('refresh', function(event) {
        event.onComplete = function () {
            var x = document.getElementById('mojoGroupFilter');
            if (x != null) {
                x.value = app.groupFilter;
            }
        }
    });


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
    //          Add Group Form
    //------------------------------------------------------------------------
    $().w2form({
        name: 'addGroupForm',
        style: 'border: 0px; background-color: #e0e0fb;',
        // url: '/v1/person',
        formURL: '/html/formAddGroupMembership.html',
        fields: [
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
                    },
                },
            },
        ]
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

    $().w2layout({
        name: 'personFormAndAddGroup',
        padding: 0,
        panels: [
            { type: 'left',    size: 0,     hidden: true,   content: 'left'    },
            { type: 'top',     size: 0,     hidden: false,  content: 'top',     resizable: false, style: app.pstyle },
            { type: 'main',    size: '60%', hidden: false,  content: 'main',    resizable: true, style: app.pstyle }, // personForm
            { type: 'preview', size: 0,     hidden: true,   content: 'preview', resizable: true,  style: app.pstyle },
            { type: 'bottom',  size: '60px',hidden: false,  content: 'bottom',  resizable: false, style: app.pstyle }, // form action buttons
            { type: 'right',   size: 0,     hidden: true,   content: 'right',   resizable: true,  style: app.pstyle }
        ]
    });

    //------------------------------------------------------------------------
    // personFormLayout
    //------------------------------------------------------------------------
    $().w2layout({
        name: 'personFormLayout',
        padding: 0,
        panels: [
            { type: 'left',    size: 0,     hidden: true,  content: 'left'    },
            { type: 'top',     size: '70%', hidden: true,  content: 'top',     resizable: true,  style: app.pstyle }, // personFormAndAddGroup
            { type: 'main',    size: '275px',hidden: false, content: 'main',   resizable: true, style: app.pstyle }, // personGroupsGrid
            { type: 'preview', size: 0,     hidden: true,  content: 'preview', resizable: true,  style: app.pstyle },
            { type: 'bottom',  size: '65px',hidden: false, content: 'bottom',  resizable: false, style: app.pstyle }, // form action buttons
            { type: 'right',   size: 0,     hidden: true,  content: 'right',   resizable: true,  style: app.pstyle }
        ]
    });
}

function resetFormExposure() {
    w2ui.personFormAndAddGroup.hide('top');
    w2ui.personFormAndAddGroup.show('main');
    w2ui.personFormAndAddGroup.hide('preview');
    w2ui.personFormAndAddGroup.show('bottom');

    w2ui.personFormLayout.show('top');
    w2ui.personFormLayout.show('main');
    w2ui.personFormLayout.hide('preview');
    w2ui.personFormLayout.show('bottom');
}
//---------------------------------------------------------------------------------
// finishPersonForm - load the layout properly.
//---------------------------------------------------------------------------------
function finishPersonForm() {
    resetFormExposure();
    w2ui.personFormAndAddGroup.content('main',w2ui.personForm);
    w2ui.personFormAndAddGroup.content('bottom',w2ui.addGroupForm);

    w2ui.personFormLayout.content('top',w2ui.personFormAndAddGroup );
    w2ui.personFormLayout.content('main',w2ui.personGroupsGrid );
    w2ui.personFormLayout.content('bottom', w2ui.personFormBtns);
}

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
    if ( typeof w2ui.addGroupForm.record === "undefined" ) {return;}
    if ( typeof w2ui.addGroupForm.record.GroupName === "undefined" ) {return;}
    if ( w2ui.addGroupForm.record.GroupName === null) {return;}

    var g = w2ui.personGroupsGrid;
    var r = g.records;

    if (addGroup.set) {
        // if we already have this group, then just return
        for (var i = 0; i < r.length; i++) {
            if ( r[i].GID === addGroup.GID) {
                w2ui.addGroupForm.record.GroupName = {};
                w2ui.addGroupForm.refresh();
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
    w2ui.addGroupForm.record.GroupName = {};
    w2ui.addGroupForm.refresh();
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
