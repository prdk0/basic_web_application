function modalOpen(csrfToken, roomId) {
    document.getElementById("check-availability-button").addEventListener('click', () => {
                let html = `
                    <form id="check-availability-form" action="" method="post" class="needs-validation" novalidate>
                    <div class="row" id="date_picker">
                        <div class="col">
                            <input required class="form-control" type="text" name="start" id="start" placeholder="arrival" required disabled>
                        </div>
                        <div class="col">
                            <input required class="form-control" type="text" name="end" id="end" placeholder="departure"  required disabled>
                        </div>
                    </div>
                    </form>
                `;
                    alertPrompt.custom({
                    msg:html, 
                    title:"choose your dates.",
                    willOpen: () => {
                        let calElement = document.getElementById("date_picker")
                        const rangepicker = new DateRangePicker(calElement, {
                            buttonClass: 'btn',
                            format: "yyyy-mm-dd",
                            minDate: "Date",
                            autohide: true
                        });
                        },
                    didOpen: () => {
                        document.getElementById("start").removeAttribute('disabled'),
                        document.getElementById("end").removeAttribute('disabled')
                    },
                    callback: function(result) {
                    let form = document.getElementById("check-availability-form");
                    let formData = new FormData(form);
                    formData.append("csrf_token", csrfToken)
                    formData.append("room_id", roomId)
                        fetch("/search-availability-json", {
                            method: "post",
                            body:formData,
                        })
                            .then(response => response.json())
                            .then(data => {
                                if(data.ok) {
                                    alertPrompt.custom({
                                        icon: "success",
                                        showConfirmButton: false,
                                        msg: '<p>Room is available!</p>'
                                              + '<p><a href="/book-room?id='
                                              + data.room_id
                                              + '&s='
                                              + data.start_date
                                              +'&e='
                                              + data.end_date
                                              +'" class="btn btn-primary">' 
                                              + 'Book Now!</a></p>'
                                        
                                    })
                                } else {
                                    alertPrompt.error({msg: "No availability!"})
                                }
                            })
                    }
                    })
                });
}