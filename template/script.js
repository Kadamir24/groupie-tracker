// const input = document.getElementById("search-input");
// const searchBtn = document.getElementById("search-btn");

// const expand = () => {
//     searchBtn.classList.toggle("close");
//     input.classList.toggle("square");
// };

// searchBtn.addEventListener("click", expand);




var res = info.city_name.filter(function(item){
    return item.city_name =='City A';
  });    
  
  //add this
  res[0].members_names = res[0].members.reduce( (acc, curr) => {
      acc += ` ${curr.members_name} `
      return acc;
  }, '');
  