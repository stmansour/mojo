"use strict";
"esversion: 8";
/*global
    w2ui, $, app, console, w2utils, document,
    getStatusString, showForm, initializePersonRecord,openGroupForm,
    setInnerHTML,

*/



// initializePropertyRecord returns an empty group record
//------------------------------------------------------------------------------
function initializePropertyRecord() {
    var rec = {
        recid: 0,
        GID: 0,
        GroupName: "",
        GroupDescription: "",
        MemberCount: 0,
        MailToCount: 0,
        OptOutCount: 0,
        BouncedCount: 0,
        ComplaintCount: 0,
        SuppressedCount: 0,
        LastModTime: null,
        LastModBy: null,
    };


    return rec;
}

function createGroupUI() {

    //------------------------------------------------------------------------
    //          groupsGrid
    //------------------------------------------------------------------------
    $().w2grid({
        name: 'groupsGrid',
        url: '/v1/groups',
        show: {
            header:       false,
            toolbar:      true,
            toolbarAdd:   true,    // indicates if toolbar add new button is visible
            footer:       true,
            lineNumbers:  false,
            selectColumn: false,
            expandColumn: false
        },
        columns: [
            {field: 'recid',            caption: 'recid',             size: '10px',  sortable: false, hidden: true},
            {field: 'GID',              caption: 'GID',               size: '100px', sortable: true, hidden: true},
            {field: 'GroupName',        caption: 'Group Name',        size: '100px', sortable: true, hidden: false},
            {field: 'GroupDescription', caption: 'Group Description', size: '300px', sortable: true, hidden: false},
            {field: 'LastModTime',      caption: 'LastModTime',       size: '100px', sortable: true, hidden: true},
        ],
        onAdd: function (/*event*/) {
            openGroupForm(0);
        },
        onClick: function(event) {
            // var sel = this.getSelection(true);
            event.onComplete = function() {
                var sel = w2ui.groupsGrid.getSelection();
                if (sel.length > 0) {
                    var ob = w2ui.groupsGrid.get(sel[0]);
                    if (ob === null ) { return; }
                    openGroupForm(ob.GID);
                    /*
                    w2ui.toplayout.show('right',true);
                    w2ui.toplayout.load('right', '/html/groupdetail.html', null, function() {
                        console.log('>>>>>>>>>> get data for Group: ' + ob.GID);
                        var gid = ob.GID;
                        $.get('/v1/groupstats/' + gid)
                        .done(function(data) {
                            if (typeof data == 'string') {  // it's weird, a successful data add gets parsed as an object, an error message does not
                                var msg = JSON.parse(data);
                                if (msg.status != "success") {
                                    console.log('Response to groupcount: ' + msg.status);
                                    return;
                                }
                                document.getElementById('mojoGroupGID').innerHTML = '' + msg.record.GroupName;
                                document.getElementById('mojoLastScrapeStart').innerHTML = '' + msg.record.LastScrapeStart;
                                document.getElementById('mojoLastScrapeStop').innerHTML = '' + msg.record.LastScrapeStop;
                                document.getElementById('mojoMemberCount').innerHTML = '' + msg.record.MemberCount;
                                document.getElementById('mojoMailToCount').innerHTML = '' + msg.record.MailToCount;
                                document.getElementById('mojoOptOutCount').innerHTML = '' + msg.record.OptOutCount;
                                document.getElementById('mojoBouncedCount').innerHTML = '' + msg.record.BouncedCount;
                                document.getElementById('mojoComplaintCount').innerHTML = '' + msg.record.ComplaintCount;
                                document.getElementById('mojoSuppressedCount').innerHTML = '' + msg.record.SuppressedCount;
                            }
                        })
                        .fail(function(data) {
                            console.log('data = ' + data);
                        });
                    });
                    */
                } // else {
                   // console.log('*** NO SELECTION FOUND ***   So, select index 0');
                    //this.select(1);
                //}
            };
        },
    });



    //------------------------------------------------------------------------
    //          Group Form
    //------------------------------------------------------------------------
    $().w2form({
        name: 'groupForm',
        formURL: '/html/groupdetail.html',
        fields: [
            { field: 'recid',           type: 'int',    required: false },
            { field: 'GID',             type: 'int',    required: false },
            { field: 'GroupName',       type: 'text',   required: false },
            { field: 'GroupDescription',type: 'text',   required: false },
            { field: 'MemberCount',     type: 'hidden', required: false },
            { field: 'MailToCount',     type: 'hidden', required:  true },
            { field: 'OptOutCount',     type: 'hidden', required: false },
            { field: 'BouncedCount',    type: 'hidden', required: false },
            { field: 'ComplaintCount',  type: 'hidden', required: false },
            { field: 'SuppressedCount', type: 'hidden', required: false },
            { field: 'LastModTime',     type: 'hidden', required: false },
            { field: 'LastModBy',       type: 'hidden', required: false },
        ],
        toolbar: {
            items: [
                { id: 'btnNotes',       type: 'button', icon: 'fa fa-sticky-note-o' },
                { id: 'bt3',            type: 'spacer' },
                { id: 'btnClose',       type: 'button', icon: 'fa fa-times' },
            ],
            onClick: function (event) {
                switch(event.target) {
                case 'btnClose':
                    closeGroupForm();
                    break;
                }
            },
        },

        actions: {
            save: function () {
                var f = w2ui.groupForm;
                var rec = {
                    GroupName: f.record.GroupName,
                    GroupDescription: f.record.GroupDescription,
                };
                var params = {
                    cmd: "save",
                    record: rec,
                };
                var dat = JSON.stringify(params);
                var url = "/v1/group/" + f.record.GID;
                $.post(url, dat, null, "json")
                .done(function (data) {
                    if (data.status === 'error') {
                        f.error('ERROR: ' + data.message);
                        return;
                    }
                    w2ui.groupsGrid.reload();
                })
                .fail( function() {
                        w2ui.rentablesGrid.error('post failed to ' + url);
                        return;
                    }
                );
                closeGroupForm();
            },
        },

        onRefresh: function(/*event*/) {
            displayGroupStats();
        },
    });
}

function displayGroupStats() {
    var r = w2ui.groupForm.record;
    setInnerHTML('mojoGroupGID', '' + r.GID);
    setInnerHTML('mojoMemberCount', r.MemberCount);
    setInnerHTML('mojoMailToCount', r.MailToCount);
    setInnerHTML('mojoOptOutCount', r.OptOutCount);
    setInnerHTML('mojoBouncedCount', r.BouncedCount);
    setInnerHTML('mojoComplaintCount', r.ComplaintCount);
    setInnerHTML('mojoSuppressedCount', r.SuppressedCount);
}

function openGroupForm(gid) {
    w2ui.toplayout.show('right', true);
    if (gid < 1) {
        w2ui.groupForm.record = initializePropertyRecord();
        w2ui.toplayout.content('right', w2ui.groupForm);
    } else {
        $.get('/v1/groupstats/' + gid)
            .done(function(data) {
                if (typeof data == 'string') { // it's weird, a successful data add gets parsed as an object, an error message does not
                    var msg = JSON.parse(data);
                    if (msg.status != "success") {
                        console.log('Response to groupcount: ' + msg.status);
                        return;
                    }
                    var r = w2ui.groupForm.record;
                    r.GID = msg.record.GID;
                    r.GroupName = msg.record.GroupName;
                    r.GroupDescription = msg.record.GroupDescription;
                    r.MemberCount = msg.record.MemberCount;
                    r.MailToCount = msg.record.MailToCount;
                    r.OptOutCount = msg.record.OptOutCount;
                    r.BouncedCount = msg.record.BouncedCount;
                    r.ComplaintCount = msg.record.ComplaintCount;
                    r.SuppressedCount = msg.record.SuppressedCount;
                    displayGroupStats();
                    w2ui.toplayout.content('right', w2ui.groupForm);
                }
            })
            .fail(function(data) {
                console.log('data = ' + data);
            });
    }

}
function closeGroupForm() {
    w2ui.toplayout.hide('right',true);
}
