var bendata = [
{
  "flag": "Submitted to hackathons with earlier start date",
  "severity": "0"
},
{
  "flag": "Started before hackathon",
  "severity": "0"
},
  {
    "flag": "Started befo hackathon",
    "severity": "1"
  }, 
  {
    "flag": "Started befo hackathon",
    "severity": "0"
  }, 
  {
    "flag": "Started befo hackathon",
    "severity": "0"
  }, 
  {
    "flag": "Started befo hackathon",
    "severity": "1"
  }, 

    ]

insertFlags(bendata)

function insertFlags(data) {
  const hackauth = document.getElementById('software-header');

  for (var i = 0; i < data.length; i++) {
    var flag = `<div class = "flag" id="flag-${i}">
    <div class="flagContent">${data[i].flag}</div>
  </div>`
    hackauth.insertAdjacentHTML('beforeend', flag);
    if (data[i].severity == "1") {
      document.getElementById("flag-" + i.toString()).style.backgroundColor = "yellow";
    }
    else {
      document.getElementById("flag-" + i.toString()).style.backgroundColor = "#fffa99";
    }
  }
}