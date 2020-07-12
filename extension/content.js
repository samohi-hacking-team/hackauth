

var bendata = [
{
  "flag": "Submitted to hackathons with earlier start date",
  "severity": "0"
},
{
  "flag": "Started before hackathon",
  "severity": "1"
}, 
    ]

insertFlags(bendata)

function insertFlags(data) {
  var hackauth = document.getElementById('software-nav');

  for (var i = 0; i < data.length; i++) {
    var flag = `<div class = "flag">
    <div>${data[i].flag}</div>
    <div>${data[i].severity}</div>
  </div>`
    hackauth.insertAdjacentHTML('afterend', flag);
  }
}