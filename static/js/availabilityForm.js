function availabilityForm(input_data) {
    const {room_id, csrf_token} = input_data
    document.getElementById("check-availability-button").addEventListener("click", function(){
    let html = `
    <form id="check-availability-form" action="" method="POST" novalidate class="needs-validation">
      <div class="row">
        <div class="col">
          <div class="row" id="reservation-dates-modal">
            <div class="col">
              <input disabled required class="form-control" type="text" name="start" id="start" placeholder="Arrival">
            </div>
            <div class="col">
              <input disabled required class="form-control" type="text" name="end" id="end" placeholder="Arrival">
            </div>
          </div>
        </div>
      </div>
    </form>
    `
    attention.custom({
      msg: html,
      title: "Choose your dates",
      willOpen: () => {
        const innerElement = document.getElementById("reservation-dates-modal")
        const rp = new DateRangePicker(innerElement, {
          format: 'yyyy-mm-dd',
          showOnFocus: true,
          minDate: new Date(),
        })
      },
      didOpen: () => {
        document.getElementById('start').removeAttribute('disabled')
        document.getElementById('end').removeAttribute('disabled')
      },
      callback: (result) => {
        console.log("called")

        let form = document.getElementById('check-availability-form')
        let formData = new FormData(form)
        formData.append("csrf_token", csrf_token)
        formData.append("room_id", room_id)
        fetch('/search-availability-json', {
          method: "post",
          body: formData,
        })
          .then(response => response.json())
          .then(data => {
            if(data.ok) {
              attention.custom({
                icon: 'success',
                msg: `<p>Room is available!</p><p><a href="/book-room?id=${data.room_id}&s=${data.start_date}&e=${data.end_date}" class="btn btn-primary">Book now!</a></p>`,
                showConfirmButton: false,
              })
            } else {
              attention.error({
                msg: "No avalibilty",
              })
            }
          })
      }
    })
  })
}