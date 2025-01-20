<h1 align="center" id="title">Nexon</h1>

<p align="center"><img src="https://socialify.git.ci/Lucifer13072/Nexon/image?font=Inter&amp;forks=1&amp;issues=1&amp;language=1&amp;name=1&amp;owner=1&amp;pattern=Brick+Wall&amp;pulls=1&amp;stargazers=1&amp;theme=Dark" alt="project-image"></p>

Nexon is a modern website engine built with powerful and flexible technologies like Go, HTML, CSS, and JavaScript. The project features a user-friendly admin panel and a functional main site designed for content management and user interaction.  

## Key Features  
### Admin Panel  
- View and analyze site statistics.  
- Manage users: add, edit, and delete accounts.  
- Manage news and other content.  

### Main Website  
- Display dynamic content for users.  
- Scalable and open for adding new features.  

## Future Plans  
- Develop a website generator module, similar to Wix, allowing users to create and edit web pages through an intuitive visual interface.  

## Technologies Used  
- **Go**: Backend for handling requests and managing data.  
- **HTML/CSS**: Structure and styling of web pages.  
- **JavaScript**: Implementation of interactive elements and dynamic functionality.  

#Project Structure

project/
├── admin/
│   ├── front/                    # Admin panel HTML files
│   │   ├── login.html            # Admin login page
│   │   └── dashboard.html        # Admin dashboard page
│   ├── scripts/                  # Scripts for admin panel
├── components/
│   ├── addition_components/      # Custom scripts
│   └── scripts/                  # Embedded website scripts
├── templates/
│   └── pages/                    # Main site HTML files
├── setup.go                      # Set start settings for engine
└── main.go                       # Main Go file
