  git clone https://github.com/chengwuloo/presstest.git

  echo "" >> README.md
  
  git init
  
  git add README.md
  
  git commit -m "first commit"
  
  git remote add origin https://github.com/chengwuloo/presstest.git
  
  git push -u origin master

  git remote add origin https://github.com/chengwuloo/presstest.git
  
  git push -u origin master

  git fetch https://github.com/chengwuloo/presstest.git
  
  git pull --allow-unrelated-histories

  git lfs install
  git-lfs version
  git lfs track "*.zip"

  git add .gitattributes
  git commit -am "modify"
  git push

  git add src/thirdpart.zip
  git commit -am "modify"
  git push
  git push origin master

  git rm --cached 
  git rm --cached -r .
  git commit --amend
  git push
