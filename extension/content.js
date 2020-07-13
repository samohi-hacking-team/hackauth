var bendata = [
{
  "flag": "Submitted to hackathons with earlier start date",
  "severity": "0"
},
{
  "flag": "Started before hackathon",
  "severity": "1"
},
  {
    "flag": "Started before hackathon",
    "severity": "1"
  }, 

    ]

insertFlags(bendata)

function insertFlags(data) {
  var hackauth = document.getElementById('software-header');

  for (var i = 0; i < data.length; i++) {
    var flag = `<div class = "flag">
    <div class="flagContent">${data[i].flag}</div>
    <div class="flagContent">${data[i].severity}</div>
  </div>`
    hackauth.insertAdjacentHTML('afterend', flag);
  }
}