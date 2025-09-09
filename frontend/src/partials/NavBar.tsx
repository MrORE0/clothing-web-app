import "./partials-styles/additional-styling.css"
import { Link } from "react-router"
import { useRef, useState, useEffect, type JSX } from "react"

//images imports 
import male_model_nav_menu from "../assets/male_model_nav_menu.jpg"
import female_model_nav_menu from "../assets/female_model_nav_menu.jpg"


import {ManClothingTypes, ManShoeTypes, ManAccessoriesTypes, WomanClothingTypes, WomanShoeTypes, WomanAccessoriesTypes, Brands} from "../additional/navExpandMenuTypes"

export const NavBar = () => {

    //State variables for all lists of the nav links
    const[clothesList, setClothesList] = useState<string>("");
    const[shoesList, setShoesList] = useState<string>("");
    const[accessoriesList, setAccessoriesList] = useState<string>("");
    const[brandsList, setBrandsList] = useState<string>("");

    // Init scroll value 
    const lastScrollY = useRef(0);

    //Navbar refference
    const refNavbar = useRef<HTMLDivElement | null>(null)

    //Pretty self-explanatory
    const refCollapseNav = useRef<HTMLDivElement | null>(null);

    //Check if scrolled up or down
    const scrollCheck = () => {
        const currentScrollY = window.scrollY;
        
        if(currentScrollY > lastScrollY.current){
            NavBarOnScrollDown();
        }else if(currentScrollY < lastScrollY.current){
            NavBarOnScrollUp();
        }
        lastScrollY.current = currentScrollY;
    }

    const NavBarOnScrollDown = () => {
        if(refNavbar.current){
            refNavbar.current.style.height = '0px';
            refNavbar.current.style.opacity = "0";
        }
    }

    const NavBarOnScrollUp = () => {
        if(refNavbar.current){
            refNavbar.current.style.height = '110px';
            refNavbar.current.style.opacity = "1";
        }
    }

    useEffect(() => {
        //listen for scrolling event
        window.addEventListener("scroll", scrollCheck)

        //Setting the brands list here since it doesn't require change later
        setBrandsList(Brands)
    }, [])

    return (
    <>
    <nav className="navbar navbar-expand-lg navbar-dark gradient-custom border-bottom border-dark fixed-top" style={{"height": "110px"}} ref={refNavbar}>
    <div className="container" >
        <Link className="nav-link home-link text-dark d-flex align-items-center fw-bold" to="/"
        style={{"fontSize": "30px", 'padding': "0px 10px"}}>Browse</Link>

        <button className="navbar-toggler" type="button" data-bs-toggle="collapse"
        data-bs-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false"
        aria-label="Toggle navigation">
        <i className="fa-solid fa-bars" style={{"color": "#000000ff"}}></i>
        </button>

        <div className="collapse navbar-collapse" id="navbarSupportedContent" ref={refCollapseNav}>
        <ul className="navbar-nav me-auto d-flex flex-row mt-3 mt-lg-0">
            <li className="nav-item dropdown position-static px-2 px-lg-1" 
            onMouseEnter={() => {    
                setClothesList(ManClothingTypes);
                setShoesList(ManShoeTypes);
                setAccessoriesList(ManAccessoriesTypes);
            }} 
            >
                <a className="nav-link active text-dark d-flex align-items-center" 
                aria-current="page" 
                href="#!"
                id="navbarDropdown"
                role="button"
                aria-expanded="false"
                >
                    Him
                </a>
                <div className="dropdown-menu w-100 dropdown-full">
                    <div className="d-flex flex-row align-items-center justify-content-around">
                        <div className="hover-content-man w-100">
                            <div className="item-types-lists-container w-100 p-5 justify-content-around d-flex flex-row">
                                <div className="clothes-list-container d-flex flex-column">
                                    <h4 className="fw-bold pb-3">Clothes</h4>
                                    <ul dangerouslySetInnerHTML={{"__html": clothesList}} className="list-unstyled"></ul>
                                </div>
                                <div className="shoes-list-container d-flex flex-column">
                                    <h4 className="fw-bold pb-3">Shoes</h4>
                                    <ul dangerouslySetInnerHTML={{"__html": shoesList}} className="list-unstyled"></ul>
                                </div>
                                <div className="accessories-list-container d-flex flex-column">
                                    <h4 className="fw-bold pb-3">Accessories</h4>
                                    <ul dangerouslySetInnerHTML={{"__html": accessoriesList}} className="list-unstyled"></ul>
                                </div>
                                <div className="brands-list-container d-flex flex-column">
                                    <h4 className="fw-bold pb-3">Brands</h4>
                                    <ul dangerouslySetInnerHTML={{"__html": brandsList}} className="list-unstyled"></ul>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </li>
            <li className="nav-item dropdown position-static text-center px-2 px-lg-1" onMouseEnter={() => {
                setClothesList(WomanClothingTypes);
                setShoesList(WomanShoeTypes);
                setAccessoriesList(WomanAccessoriesTypes);
            }}>
                <a 
                    className="nav-link active text-dark d-flex flex-column justify-content-center" 
                    aria-current="page" 
                    href="#!"
                    id="navbarDropdown"
                    role="button"
                    aria-expanded="false"
                >
                    Her
                </a>
                <div className="dropdown-menu w-100 dropdown-full">
                    <div className="d-flex flex-row align-items-center justify-content-around">
                        <div className="hover-content-man w-100">
                            <div className="item-types-lists-container w-100 p-5 justify-content-around d-flex flex-row">
                                <div className="clothes-list-container d-flex flex-column">
                                    <h4 className="fw-bold pb-3">Clothes</h4>
                                    <ul dangerouslySetInnerHTML={{"__html": clothesList}} className="list-unstyled"></ul>
                                </div>
                                <div className="shoes-list-container d-flex flex-column">
                                    <h4 className="fw-bold pb-3">Shoes</h4>
                                    <ul dangerouslySetInnerHTML={{"__html": shoesList}} className="list-unstyled"></ul>
                                </div>
                                <div className="accessories-list-container d-flex flex-column">
                                    <h4 className="fw-bold pb-3">Accessories</h4>
                                    <ul dangerouslySetInnerHTML={{"__html": accessoriesList}} className="list-unstyled"></ul>
                                </div>
                                <div className="brands-list-container d-flex flex-column">
                                    <h4 className="fw-bold pb-3">Brands</h4>
                                    <ul dangerouslySetInnerHTML={{"__html": brandsList}} className="list-unstyled"></ul>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </li>
            <li className="nav-item dropdown position-static text-center px-2 px-lg-1" onMouseEnter={() => {
                setClothesList(WomanClothingTypes);
                setShoesList(WomanShoeTypes);
                setAccessoriesList(WomanAccessoriesTypes);
            }}>
            <a className="nav-link active text-dark d-flex flex-column justify-content-center" href="#!">
                Hot
            </a>
            <div className="dropdown-menu w-100 dropdown-full">
                    <div className="d-flex flex-row align-items-center justify-content-around">
                       <div className="card border-0 nav-card-item d-flex align-items-center" >
                            <img className="card-img-top" src={male_model_nav_menu} alt="Male Model Image Navbar"/>
                            <h3 className="card-body fw-bold">For Him</h3>
                        </div>
                        <div className="card border-0 nav-card-item d-flex align-items-center">
                            <img className="card-img-top" src={female_model_nav_menu} alt="Female Model Image Navbar"/>
                            <h3 className="card-body fw-bold">For Her</h3>
                        </div>
                    </div>
                </div>
            </li>
            <li className="nav-item text-center px-2 px-lg-1">
            <Link to="/customize" className="nav-link active text-dark d-flex flex-column justify-content-center">
                Custome
            </Link>
            </li>

        </ul>

        <div className="brand-container">
            <h1 className="fw-bold" style={{"fontFamily": "'Bright Live', sans-serif", "fontSize": "60px"}}>AssK</h1>
        </div>


        <ul className="navbar-nav ms-auto d-flex flex-row mt-3 mt-lg-0">
            <li className="nav-item text-center px-2 px-lg-1">
            <a className="nav-link text-dark d-flex flex-column justify-content-center" href="#!">
                <div>
                <i className="fa-solid fa-circle-user fa-2xl" style={{"color": "#280028"}}></i>
                </div>
            </a>
            </li>
        </ul>

        <form className="d-flex input-group w-auto ms-lg-3 my-3 my-lg-0">
            <input type="search" className="form-control" placeholder="Explore" aria-label="Explore" />
            <button type="button" className="btn btn-outline-dark" data-mdb-ripple-init data-mdb-ripple-color="dark">
            Find
            </button>
        </form>
        </div>
    </div>
    </nav>
    </>
    )
}