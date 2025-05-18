import "./Footer.css"
import { useState } from "react"

function Footer() {

  const [updateTime, setUpdateTime] = useState(new Date())

  let updateTimeString = "Unknow"

  /* Update time will need to be gathered from the server */
  
  return (
    <div className="footer_main">
      <p>Last Updated: {updateTimeString}</p>
      <div className="credits">
        <p>Made By: <a>Nathaniel Martes</a> and <a>Zackery Drake</a>
        </p>
      </div>
    </div>
  );
}

export default Footer;
