import "./Header.css"
import search_icon from "../assets/search-icon.svg"
import filter_icon from "../assets/filter-icon.svg"
import bars_icon from "../assets/bars-solid.svg"
import xmark from "../assets/xmark.svg"
import Slider from '@mui/material/Slider'
import { useState, useEffect } from "react"

function Header({setFilter}) {
  const [semester, setSemester] = useState("Sp")
  const [year, setYear] = useState(0)
  const [years, setYears] = useState([])
  const [filterVisible, setFilterVisible] = useState(false)
  const [value, setValue] = useState([8, 22])
  const [otherSelected, isOther] = useState(false)
  const [startTime, setStartTime] = useState(8)
  const [endTime, setEndTime] = useState(22)

  const handleSliderChange = (event, newValue) => {
    setValue(newValue);
    setStartTime(newValue[0])
    setEndTime(newValue[1])
  };

  const handleRadioChange = (event) => {
    isOther(event.target.value === "other")
  }

  const search = (searchValue) => {
    setFilter(searchValue);
  }

  const filter = (filterForm) => {
    // Set filter parameters here
    // TODO: pull parameters from form and pass into setFilter
    const parameters = {}
    setFilter(parameters);
  }

  const validateForm = (event) => {
    event.preventDefault();
    const events = {
      "searchForm": search(document.forms["searchForm"]["searchBox"].value),
      "filterForm": filter(document.forms["filterForm"])
    };
    events[event.target.name];
  }

  useEffect( () => {
    /* Simulate Data */
    const tmpYears = []
    for (let i = 25; i >= 18; i--) {
      tmpYears.push(i)
    }
    setYears(tmpYears)
    /* ASSUME years always eixsts */
    setYear(years[0])
  }, [])

  function toggleFilter() {
    setFilterVisible(!filterVisible)
  }

  function closeFilter() {
    setFilterVisible(false)
  }

  return (
    <>
      <div className="header_main">
        <nav className="pages">
          <a>Courses</a>
          <a>My Courses</a>
        </nav>
        <div className="search_and_filter">
          <button className="filter_button" onClick={toggleFilter}>
            <img className="icon" src={filter_icon}/>
          </button>
          <select value={semester} onChange={(event) => setSemester(event.target.value)}>
            <option value="F">Fall</option>
            <option value="Sp">Spring</option>
          </select>
          <select>
            {years.map((year, index) => (<option value={year}>{year}</option>))}
          </select>
          <form className="search_and_filter" name="searchForm" onSubmit={validateForm}>
            <input type="search" name="searchBox" placeholder="Search Courses"/>
          </form>
          <button className="search_button">
            <img className="icon" src={search_icon}/>
          </button>
        </div>
      </div>

      <form className="filter_panel" name="filterForm" onSubmit={validateForm}>
        <div className={`${filterVisible ? "visible" : ""} filter_main`}>
          <button className="close_button" onClick={closeFilter}>
            <img className="icon" src={xmark}/>
          </button>
          <div>
            <p className="filter_option">Delivery Mode</p>
            <input type="checkbox" className="check-opt" name="f2f-opt" id="f2f-opt"/>
            <label htmlFor="f2f-opt" className="opt-label">F2F</label>
            <br></br>
            <input type="checkbox" className="check-opt" name="hyb-opt" id="hyb-opt"/>
            <label htmlFor="hyb-opt" className="opt-label">HYB</label>
            <br></br>
            <input type="checkbox" className="check-opt" name="ol-opt" id="ol-opt"/>
            <label htmlFor="ol-opt" className="opt-label">OL</label>
            <br></br>
          </div>
          <div>
            <p className="filter_option">Course Category</p>
            <textarea id="instructor-opt" rows="3" cols="50" placeholder="MTH, ENG, etc." spellCheck="false" name="category_text_box"></textarea>
          </div>
          <div>
            <p className="filter_option">Credits</p>
            <select id="credits-opt" className="dropdown-opt" name="credits-dropdown">
              <option value=""></option>
              <option value="0">0</option>
              <option value="1">1</option>
              <option value="2">2</option>
              <option value="3">3</option>
              <option value="4">4</option>
              <option value="5">5</option>
              <option value="6">6</option>
            </select>
          </div>
          <div>
            <p className="filter_option">Day</p>
            <input type="checkbox" className="check-opt" name="m-opt" id="m-opt"/>
            <label htmlFor="m-opt" className="opt-label">M</label>
            <br></br>
            <input type="checkbox" className="check-opt" name="t-opt" id="t-opt"/>
            <label htmlFor="t-opt" className="opt-label">T</label>
            <br></br>
            <input type="checkbox" className="check-opt" name="w-opt" id="w-opt"/>
            <label htmlFor="w-opt" className="opt-label">W</label>
            <br></br>
            <input type="checkbox" className="check-opt" name="r-opt" id="r-opt"/>
            <label htmlFor="r-opt" className="opt-label">R</label>
            <br></br>
            <input type="checkbox" className="check-opt" name="f-opt" id="f-opt"/>
            <label htmlFor="f-opt" className="opt-label">F</label>
            <br></br>
          </div>
          <div>
            <p className="filter_option">Time</p>
            <input type="radio"
                   name="time-am"
                   className="check-opt"
                   id="am-opt"
                   value="am"
                   onChange={handleRadioChange}
            />
            <label htmlFor="am-opt" className="opt-label">AM</label>
            <input type="radio"
                   name="time-pm"
                   className="check-opt"
                   id="pm-opt"
                   value="pm"
                   onChange={handleRadioChange}
            />
            <label htmlFor="pm-opt" className="opt-label">PM</label>
            <input type="radio"
                   name="time-other"
                   className="check-opt"
                   id="other-opt"
                   value="other"
                   onChange={handleRadioChange}
            />
            <label htmlFor="other-opt" className="opt-label">Other</label>
            <br></br>
            <br></br>
            <Slider id="slider"
                    name="time-slider"
                    value={value}
                    onChange={handleSliderChange}
                    disabled={!otherSelected}
                    disableSwap
                    min={8}
                    max={22}
            />
            <p id="time">{startTime < 13  ? startTime : startTime%13 + 1}{startTime/12 < 1 ? "am" : "pm"}
              &nbsp;-&nbsp;
            {endTime < 13 ? endTime : endTime%13 + 1}{endTime/13 < 1 ? "am" : "pm"}</p>
          </div>
          <div>
            <p className="filter_option">Location</p>
            <select id="location-opt" className="dropdown-opt" name="location-dropdown">
              <option value=""></option>
              <option value="SIMMS">SIMMS</option>
              <option value="BEDFRD">BEDFRD</option>
              <option value="BEDFRD ANNEX">BEDFRD ANNEX</option>
              <option value="BREIS">BREIS</option>
              <option value="FARLEY">FARLEY</option>
              <option value="MARTS">MARTS</option>
              <option value="SLC">SLC</option>
              <option value="CSC">CSC</option>
              <option value="CONHAM">CONHAM</option>
              <option value="KARAM">KARAM</option>
              <option value="DDD">DDD</option>
              <option value="KIRBY">KIRBY</option>
            </select>
          </div>
          <div>
            <p className="filter_option">Instructor</p>
            <textarea id="instructor-opt" rows="3" cols="50" placeholder="Instructor name" name="instructor-text-box" spellCheck="false"></textarea>
          </div>
          <div>
            <p className="filter_option">Misc</p>
            <input type="checkbox" className="check-opt" name="hide-closed" id="hide-closed-opt"/>
            <label htmlFor="hide-closed-opt" className="opt-label">Hide Closed</label>
            <br></br>
            <br></br>
            <label htmlFor="course-quant-opt" id="per_page_label" className="opt-label">Courses per page</label>
            <select id="course-quant-opt" className="dropdown-opt" name="display-num">
              <option value="0">All</option>
              <option value="5">0</option>
              <option value="10">10</option>
              <option value="15">15</option>
              <option value="20">20</option>
              <option value="50">50</option>
              <option value="100">100</option>
            </select>
          </div>
          <br></br>
          <br></br>
          <br></br>
          <br></br>
          <button className="filter_action" onClick={closeFilter} name="apply-button">
            Apply
          </button>
        </div>
      </form>
    </>
  );
}

export default Header;
