function CheckRoomAvailability(room_id, csrf_token){
    document.getElementById("check-availability-button").addEventListener("click", function () {
        let html = `
    <form id="check-availability-form" action="" method="post" novalidate class="needs-validation">
        <div class="form-row">
            <div class="col">
                <div class="form-row" id="reservation-dates-modal">
                    <div class="col">
                        <input disabled required class="form-control" type="text" name="start" id="start" placeholder="Arrival">
                    </div>
                    <div class="col">
                        <input disabled required class="form-control" type="text" name="end" id="end" placeholder="Departure">
                    </div>

                </div>
            </div>
        </div>
    </form>
    `;
    attention.custom({
    title: "Choose your dates",
    msg: html,
    willOpen: () => {
        const elem = document.getElementById("reservation-dates-modal");
        const rp = new DateRangePicker(elem, {
        format: "yyyy-mm-dd",
        showOnFocus: true,
        minDate: new Date(),
        });
    },
    didOpen: () => {
        document.getElementById("start").removeAttribute("disabled");
        document.getElementById("end").removeAttribute("disabled");
    },
    callback: function (result) {
        console.log("called");

        let form = document.getElementById("check-availability-form");
        let formData = new FormData(form);
        formData.append("csrf_token", csrf_token);
        formData.append("room_id", room_id);

        fetch("/search-availability-json", {
        method: "post",
        body: formData,
        })
        .then((response) => response.json())
        .then((data) => {
            if (data.ok) {
            attention.custom({
                icon: "success",
                msg:
                "<p>Room is available!</p>" +
                `<p><a href="/book-room?id=${data.room_id}&s=${data.start_date}&e=${data.end_date}" 
                    class="btn btn-primary">` +
                "Book now!</a></p>",
                showConfirmButton: false,
                showCancelButton: false,
            });
            } else {
            attention.error({
                msg: "No availability",
            });
            }
        });
    },
    });
});
}