python3 -m venv venv
source venv/bin/activate
pip install -r utils/requirements.txt
python utils/pb_build.py --all .
npm install --prefix ./frontend