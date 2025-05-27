import "./Header.css"
import search_icon from "../assets/search-icon.svg"
import filter_icon from "../assets/filter-icon.svg"
import bars_icon from "../assets/bars-solid.svg"
import xmark from "../assets/xmark.svg"

import { useState, useEffect } from "react"

function Header() {
  
  const [semester, setSemester] = useState("Sp")
  const [year, setYear] = useState(0)
  const [years, setYears] = useState([])
  const [filterVisible, setFilterVisible] = useState(false)

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
        <input type="search" placeholder="Search Courses"/>
        <button className="search_button">
          <img className="icon" src={search_icon}/>
        </button>
      </div>
    </div>
    <div className={`${filterVisible ? "visible" : ""} filter_main`}>
      <button className="close_button" onClick={closeFilter}>
        <img className="icon" src={xmark}/>
      </button>

      <p>I am an empty filter, I need stuff</p>


    </div>
    </>
  );
}

export default Header;
