import "./Card.css";
import copy from "../assets/copy-link-icon.svg"
import { useState } from "react"

function Card({header = "Header", title = "Title", extra_info = "Extra Info", 
              time = "MWF 9 - 9:50AM SLC 101", crn = "00000", 
              credits = 3.00, students = 0, limit = 1}) {

  const [message, setMessage] = useState("Add")
  const buttonStyles = {
    backgroundColor: "#002855"
  }

  // const checkStudents = (value) => {
  //   const card = document.getElementById("card_button");
  //   if (value == 1) {
  //       buttonStyles.backgroundColor = "red"
  //       setMessage("Closed")
  //   }
  // };
  // checkStudents(students / limit)

  const isFull = students >= limit;
  const buttonLabel = isFull ? "Closed" : "Add";
  const buttonStyle = {
    backgroundColor: isFull ? "red" : "#002855"
  };

  return (
      <div className="main">
        <p className="header">{header}</p>
        <p className="title">{title}</p>
        <p className="credits">Credits: {parseFloat(credits).toFixed(2)}</p>
        <p className="extra_info">{extra_info}</p>
        <p className="time">{time}</p>
        <div className="bottom_bar">
          <div className="crn">
            <p>CRN: {crn}</p>
            <img src={copy} alt="Copy" width="24" height="24"/>
          </div>
          <div className="padding"></div>
          <div className="students">
            <p> {students} / {limit}</p>
          </div>
          <button className="card_button" id="card_button" style={buttonStyle}>{buttonLabel}</button>
          </div>
      </div>
  );
}

export default Card;
